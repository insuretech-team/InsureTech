"use client";

/**
 * EKYCFlow — React port of the Svelte FLVE client (UI/UX parity)
 *
 * Calls InsureTech BFF (NOT direct FLVE):
 *  - POST /api/auth/kyc/initiate   → {session_id, steps, session_state}
 *  - POST /api/auth/kyc/frame      → frame response with detection/liveness
 *  - POST /api/auth/kyc/complete   → finalises and sets PENDING_REVIEW status
 *
 * Canvas overlay (live video):
 *  - Face guide oval: 40% width, 60% height (matches Svelte)
 *  - Bounding box: pixel coords from FLVE scaled to display rect
 *  - Eye contours: edge-list drawn from FLVE FaceMesh points — called from
 *    processFrame with actual videoWidth/videoHeight (NOT from RAF)
 *  - Head pose indicator: top-right corner
 *
 * Key fix: drawOverlay(w, h) is called inside processFrame() with the real
 * captured frame dimensions — identical to Svelte. The separate RAF loop is
 * removed because it ran with stale/zero videoWidth and caused eye contours
 * to be scaled wrong or skipped entirely.
 */

import React, {
  useRef, useCallback, useState, useEffect,
} from "react";
import {
  Camera, CameraOff, CheckCircle2, XCircle,
  AlertCircle, Loader2,
} from "lucide-react";
// Card/Button no longer used — UI is custom dark-theme to match Svelte client

// ── Types ──────────────────────────────────────────────────────────────────

interface EKYCStep {
  step_number: number;
  type?: string;
  challenge_type?: string;
  state: string;
  instruction: string;
  timeout_seconds: number;
  confidence: number;
}

interface Detection {
  detected: boolean;
  box?: {
    x: number;
    y: number;
    width: number;
    height: number;
  };
}

interface HeadPose {
  yaw: number;
  pitch: number;
  roll: number;
}

interface EyeState {
  left_openness: number;
  right_openness: number;
  is_blinking: boolean;
}

interface EyeContour {
  edges: [number, number][];
  points: Record<string, { x: number; y: number }>;
}

interface FrameResponse {
  session_state: string;
  guidance?: string[];
  guidance_messages?: string[];
  liveness_score?: number;
  liveness_confidence?: number;
  overall_progress?: number;
  step_completed?: boolean;
  detection?: Detection;
  head_pose?: HeadPose;
  eye_state?: EyeState;
  // eye_contours_json: JSON-encoded string from Go backend (field 17 in proto).
  // The backend JSON-encodes the FLVE eye mesh to avoid typed proto messages.
  eye_contours_json?: string;
  // eye_contours: parsed from eye_contours_json (done in processFrame)
  eye_contours?: {
    left?: EyeContour;
    right?: EyeContour;
  };
  current_step?: EKYCStep;
  current_step_detail?: EKYCStep;
  error?: string;
}

interface InitiateResponse {
  ok: boolean;
  session_id?: string;
  steps?: EKYCStep[];
  session_state?: string;
  message?: string;
}

interface CompleteResponse {
  ok: boolean;
  profile_image_url?: string;
  captured_image_base64?: string;
  liveness_confidence?: number;
  message?: string;
}

// ── Challenge config ───────────────────────────────────────────────────────

const CHALLENGE_CONFIG: Record<string, {
  icon: string;
  label: string;
  description: string;
}> = {
  // Stripped (after stripPrefix)
  BLINK:      { icon: "👁️",  label: "Blink",      description: "Blink your eyes naturally" },
  LOOK_LEFT:  { icon: "👈",  label: "Look Left",   description: "Turn your head to the left" },
  LOOK_RIGHT: { icon: "👉",  label: "Look Right",  description: "Turn your head to the right" },
  CAPTURE:    { icon: "📸",  label: "Look Ahead",  description: "Look straight at the camera" },
  LIVENESS:   { icon: "📸",  label: "Look Ahead",  description: "Look straight at the camera" },
  // Full enum strings (in case BFF returns them unstripped)
  EKYC_CHALLENGE_BLINK:      { icon: "👁️",  label: "Blink",      description: "Blink your eyes naturally" },
  EKYC_CHALLENGE_LOOK_LEFT:  { icon: "👈",  label: "Look Left",   description: "Turn your head to the left" },
  EKYC_CHALLENGE_LOOK_RIGHT: { icon: "👉",  label: "Look Right",  description: "Turn your head to the right" },
  EKYC_CHALLENGE_CAPTURE:    { icon: "📸",  label: "Look Ahead",  description: "Look straight at the camera" },
};

// Default steps shown before/during session (mirrors Svelte — always visible)
const DEFAULT_STEPS: EKYCStep[] = [
  { step_number: 1, type: "EKYC_CHALLENGE_BLINK",      challenge_type: "EKYC_CHALLENGE_BLINK",      state: "EKYC_STEP_PENDING", instruction: "Blink",      timeout_seconds: 10, confidence: 0 },
  { step_number: 2, type: "EKYC_CHALLENGE_LOOK_LEFT",  challenge_type: "EKYC_CHALLENGE_LOOK_LEFT",  state: "EKYC_STEP_PENDING", instruction: "Look Left",  timeout_seconds: 10, confidence: 0 },
  { step_number: 3, type: "EKYC_CHALLENGE_LOOK_RIGHT", challenge_type: "EKYC_CHALLENGE_LOOK_RIGHT", state: "EKYC_STEP_PENDING", instruction: "Look Right", timeout_seconds: 10, confidence: 0 },
  { step_number: 4, type: "EKYC_CHALLENGE_CAPTURE",    challenge_type: "EKYC_CHALLENGE_CAPTURE",    state: "EKYC_STEP_PENDING", instruction: "Look Ahead", timeout_seconds: 10, confidence: 0 },
];

// ── Utilities ──────────────────────────────────────────────────────────────

function stripPrefix(raw?: string): string {
  if (!raw) return "";
  return raw
    .replace(/^EKYC_CHALLENGE_/, "")
    .replace(/^EKYC_SESSION_/, "")
    .replace(/^EKYC_STEP_/, "");
}

function coerceGuidance(value: any): string[] {
  if (!value) return [];
  if (Array.isArray(value)) return value;
  return [String(value)];
}

// ── Props ──────────────────────────────────────────────────────────────────

interface EKYCFlowProps {
  userId: string;
  onComplete?: (profileImageUrl: string) => void;
  onError?: (message: string) => void;
}

// ── Component ──────────────────────────────────────────────────────────────

const FRAME_INTERVAL_MS = 100;

export function EKYCFlow({ userId, onComplete, onError }: EKYCFlowProps) {
  // DOM refs
  const videoRef   = useRef<HTMLVideoElement>(null);
  const canvasRef  = useRef<HTMLCanvasElement>(null);
  const overlayRef = useRef<HTMLCanvasElement>(null);
  const streamRef  = useRef<MediaStream | null>(null);
  const timerRef   = useRef<ReturnType<typeof setTimeout> | null>(null);

  // Refs that drive the frame loop — never stale in closures (mirrors Svelte's reactivity)
  const sessionIdRef    = useRef<string>("");
  const isRunningRef    = useRef<boolean>(false);
  const sessionStateRef = useRef<string>("");

  // React state
  const [phase, setPhase] = useState<"idle" | "starting" | "active" | "completing" | "done" | "error">("idle");
  const [steps, setSteps] = useState<EKYCStep[]>([]);
  const [currentStepIdx, setCurrentStepIdx] = useState(0);
  const [overallProgress, setOverallProgress] = useState(0);
  const [guidance, setGuidance] = useState<string[]>([]);
  const [faceDetected, setFaceDetected] = useState(false);
  const [livenessScore, setLivenessScore] = useState(0);
  const [livenessConfidence, setLivenessConfidence] = useState(0);
  const [errorMsg, setErrorMsg] = useState("");
  const [inlineError, setInlineError] = useState("");
  const [profileImageUrl, setProfileImageUrl] = useState("");
  const [capturedImageB64, setCapturedImageB64] = useState("");
  const [frameCount, setFrameCount] = useState(0);
  const [avgLatency, setAvgLatency] = useState(0);
  const [lastFrameResponse, setLastFrameResponse] = useState<FrameResponse | null>(null);

  const latenciesRef = useRef<number[]>([]);

  // ── Camera control ────────────────────────────────────────────────────────

  const startCamera = useCallback(async () => {
    const stream = await navigator.mediaDevices.getUserMedia({
      video: { width: 640, height: 480, facingMode: "user" },
      audio: false,
    });
    streamRef.current = stream;
    if (videoRef.current) {
      videoRef.current.srcObject = stream;
      await videoRef.current.play();
    }
  }, []);

  const stopCamera = useCallback(() => {
    isRunningRef.current = false;  // stop frame loop immediately
    if (timerRef.current) {
      clearTimeout(timerRef.current);
      timerRef.current = null;
    }
    if (streamRef.current) {
      streamRef.current.getTracks().forEach(t => t.stop());
      streamRef.current = null;
    }
    if (overlayRef.current) {
      const ctx = overlayRef.current.getContext("2d");
      if (ctx) ctx.clearRect(0, 0, overlayRef.current.width, overlayRef.current.height);
    }
  }, []);

  // ── Overlay drawing ───────────────────────────────────────────────────────
  // Called from processFrame() with the actual captured frame dimensions —
  // identical to Svelte. Do NOT use a separate RAF loop; videoWidth/videoHeight
  // can be 0 before the video element renders its first frame which makes all
  // scale calculations wrong and eye contours invisible.

  const drawOverlay = useCallback((
    videoWidth: number,
    videoHeight: number,
    data: FrameResponse | null,
  ) => {
    const overlay = overlayRef.current;
    const video   = videoRef.current;
    if (!overlay || !video) return;

    const rect = video.getBoundingClientRect();
    overlay.width  = rect.width;
    overlay.height = rect.height;

    const ctx = overlay.getContext("2d");
    if (!ctx) return;

    ctx.clearRect(0, 0, overlay.width, overlay.height);

    const detected = data?.detection?.detected ?? false;

    // ── Face guide oval (centered, 40% width, 60% height — matches Svelte) ──
    const cX = overlay.width  / 2;
    const cY = overlay.height / 2;
    ctx.strokeStyle = detected ? "rgba(34,197,94,0.5)" : "rgba(255,255,255,0.3)";
    ctx.lineWidth = 3;
    ctx.setLineDash([10, 5]);
    ctx.beginPath();
    ctx.ellipse(cX, cY, overlay.width * 0.2, overlay.height * 0.3, 0, 0, 2 * Math.PI);
    ctx.stroke();
    ctx.setLineDash([]);

    if (!data || !videoWidth || !videoHeight) return;

    const scaleX = overlay.width  / videoWidth;
    const scaleY = overlay.height / videoHeight;

    // ── Bounding box ──
    const box = data.detection?.box;
    if (detected && box) {
      const liveness = data.liveness_score ?? data.liveness_confidence ?? 0;
      ctx.strokeStyle = liveness > 0.5 ? "#22c55e" : "#ef4444";
      ctx.lineWidth = 2;
      ctx.strokeRect(box.x * scaleX, box.y * scaleY, box.width * scaleX, box.height * scaleY);
    }

    // ── Eye contours (FaceMesh edge list from FLVE) ──
    const eyeContours = data.eye_contours;
    if (detected && eyeContours) {
      const isBlinking = data.eye_state?.is_blinking ?? false;
      const eyeColor   = isBlinking ? "rgba(34,197,94,0.95)" : "rgba(59,130,246,0.9)";

      const drawEye = (eye?: EyeContour) => {
        if (!eye?.edges?.length || !eye?.points) return;
        ctx.strokeStyle = eyeColor;
        ctx.lineWidth   = 2;
        ctx.beginPath();
        for (const [a, b] of eye.edges) {
          const pa = eye.points[String(a)];
          const pb = eye.points[String(b)];
          if (!pa || !pb) continue;
          ctx.moveTo(pa.x * scaleX, pa.y * scaleY);
          ctx.lineTo(pb.x * scaleX, pb.y * scaleY);
        }
        ctx.stroke();
      };

      drawEye(eyeContours.left);
      drawEye(eyeContours.right);
    }

    // ── Head pose indicator (top-right) ──
    const hp = data.head_pose;
    if (hp) {
      const ix = overlay.width - 80;
      const iy = 60;

      ctx.fillStyle = "rgba(0,0,0,0.5)";
      ctx.fillRect(ix - 30, iy - 30, 60, 60);

      ctx.strokeStyle = "#ffffff";
      ctx.lineWidth = 2;
      ctx.beginPath();
      ctx.arc(ix, iy, 20, 0, 2 * Math.PI);
      ctx.stroke();

      const dx = Math.sin((hp.yaw   * Math.PI) / 180) * 15;
      const dy = Math.sin((hp.pitch * Math.PI) / 180) * 15;

      ctx.fillStyle = "#22c55e";
      ctx.beginPath();
      ctx.arc(ix + dx, iy + dy, 5, 0, 2 * Math.PI);
      ctx.fill();
    }
  }, []);

  // ── eKYC API calls ────────────────────────────────────────────────────────

  const startEKYC = useCallback(async () => {
    setPhase("starting");
    setErrorMsg("");
    setInlineError("");
    setCapturedImageB64("");
    setProfileImageUrl("");
    setLivenessConfidence(0);
    setFrameCount(0);
    latenciesRef.current = [];
    setAvgLatency(0);

    try {
      await startCamera();

      const resp = await fetch("/api/auth/kyc/initiate", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ user_id: userId }),
      });

      if (!resp.ok) throw new Error(`HTTP ${resp.status}`);

      const data = await resp.json() as InitiateResponse;
      if (!data.ok || !data.session_id) {
        throw new Error(data.message ?? "Failed to start session");
      }

      sessionIdRef.current = data.session_id;
      sessionStateRef.current = "ACTIVE";
      isRunningRef.current = true;
      // Merge server steps into DEFAULT_STEPS to preserve ordering and types
      const serverSteps = data.steps ?? [];
      setSteps(serverSteps.length > 0 ? serverSteps : DEFAULT_STEPS);
      setPhase("active");
      scheduleFrame();
    } catch (e) {
      const msg = e instanceof Error ? e.message : "Failed to start eKYC";
      setErrorMsg(msg);
      setPhase("error");
      stopCamera();
      onError?.(msg);
    }
  }, [userId, startCamera, stopCamera, onError]);

  const completeSession = useCallback(async (sid: string) => {
    setPhase("completing");
    stopCamera();
    try {
      const resp = await fetch("/api/auth/kyc/complete", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ session_id: sid }),
      });

      if (!resp.ok) throw new Error(`HTTP ${resp.status}`);

      const data = await resp.json() as CompleteResponse;
      if (!data.ok) throw new Error(data.message ?? "Completion failed");

      setProfileImageUrl(data.profile_image_url ?? "");
      setCapturedImageB64(data.captured_image_base64 ?? "");
      setLivenessConfidence(data.liveness_confidence ?? 0);
      setPhase("done");
      onComplete?.(data.profile_image_url ?? "");
    } catch (e) {
      const msg = e instanceof Error ? e.message : "Completion failed";
      setErrorMsg(msg);
      setPhase("error");
      onError?.(msg);
    }
  }, [stopCamera, onComplete, onError]);

  // scheduleFrame — exact mirror of Svelte's scheduleFrame
  // Uses isRunningRef (not phase state) to avoid stale closures
  function scheduleFrame() {
    if (!isRunningRef.current) return;
    timerRef.current = setTimeout(async () => {
      try {
        await processFrame();
      } catch (e: unknown) {
        console.error("Frame error:", e);
      } finally {
        // Only reschedule when session is still ACTIVE (mirrors Svelte)
        if (isRunningRef.current && sessionStateRef.current === "EKYC_SESSION_ACTIVE") {
          scheduleFrame();
        }
      }
    }, FRAME_INTERVAL_MS);
  }

  const processFrame = async () => {
    const sid = sessionIdRef.current;
    if (!sid) return;

    const video = videoRef.current;
    const canvas = canvasRef.current;
    if (!video || !canvas) return;

    const w = video.videoWidth;
    const h = video.videoHeight;
    if (!w || !h) return;   // video not ready yet — skip frame (mirrors Svelte guard)

    canvas.width  = w;
    canvas.height = h;
    const ctx2d = canvas.getContext("2d");
    if (!ctx2d) return;
    ctx2d.drawImage(video, 0, 0, w, h);

    const blob: Blob = await new Promise((resolve, reject) =>
      canvas.toBlob(b => (b ? resolve(b) : reject(new Error("toBlob failed"))), "image/jpeg", 0.85)
    );

    const t0 = performance.now();
    try {
      const formData = new FormData();
      formData.append("user_id", userId);
      formData.append("session_id", sid);
      formData.append("file", blob, "frame.jpg");

      const resp = await fetch("/api/auth/kyc/frame", {
        method: "POST",
        body: formData,
      });

      if (!resp.ok) return;

      const raw = await resp.json() as FrameResponse;

      // Parse eye_contours_json — the Go backend JSON-encodes the FLVE eye mesh
      // into a string field to avoid adding typed proto messages for the edge list.
      const data: FrameResponse = raw;
      if (raw.eye_contours_json && !raw.eye_contours) {
        try {
          data.eye_contours = JSON.parse(raw.eye_contours_json);
        } catch {
          // ignore parse errors
        }
      }

      const latency = performance.now() - t0;
      latenciesRef.current.push(latency);
      if (latenciesRef.current.length > 30) latenciesRef.current.shift();
      setAvgLatency(latenciesRef.current.reduce((a, b) => a + b, 0) / latenciesRef.current.length);
      setFrameCount(c => c + 1);

      const rawSessionState = data.session_state ?? "";
      sessionStateRef.current = rawSessionState;

      const liveness = data.liveness_score ?? data.liveness_confidence ?? 0;
      setGuidance(coerceGuidance(data.guidance ?? data.guidance_messages ?? []));
      setLivenessScore(liveness);
      setOverallProgress(data.overall_progress ?? 0);
      setFaceDetected(data.detection?.detected ?? false);

      const currentStep = data.current_step ?? data.current_step_detail;
      if (currentStep) {
        const stepNum = currentStep.step_number ?? 1;
        setCurrentStepIdx(stepNum - 1);
        setSteps(prev => {
          const updated = [...prev];
          const idx = updated.findIndex(s => s.step_number === stepNum);
          if (idx >= 0) updated[idx] = { ...updated[idx], ...currentStep };
          // Mark previous step completed when step advances (mirrors Svelte)
          if (data.step_completed && idx > 0) {
            updated[idx - 1] = { ...updated[idx - 1], state: "EKYC_STEP_COMPLETED" };
          }
          return updated;
        });
      }

      if (data.error) setInlineError(data.error);
      else setInlineError("");

      setLastFrameResponse(data);

      // Draw overlay with REAL frame dimensions — this is the fix for eye contours
      drawOverlay(w, h, data);

      // Session state transitions
      if (rawSessionState === "EKYC_SESSION_UPLOADING" || rawSessionState === "EKYC_SESSION_COMPLETED") {
        isRunningRef.current = false;
        await completeSession(sid);
      } else if (rawSessionState === "EKYC_SESSION_FAILED" || rawSessionState === "EKYC_SESSION_EXPIRED") {
        isRunningRef.current = false;
        setErrorMsg("Session expired or failed. Please try again.");
        setPhase("error");
        stopCamera();
      }
    } catch {
      // Frame errors are non-fatal
    }
  };

  // Cleanup on unmount
  useEffect(() => () => stopCamera(), [stopCamera]);

  // ── Render helpers ─────────────────────────────────────────────────────────

  const currentStep  = steps[currentStepIdx];
  const challengeType = stripPrefix(currentStep?.type ?? currentStep?.challenge_type);
  const challengeCfg  = CHALLENGE_CONFIG[challengeType] || { icon: "🔍", label: challengeType, description: "" };

  const resetToIdle = () => {
    setPhase("idle");
    setErrorMsg("");
    setInlineError("");
    setSteps([]);
    setCapturedImageB64("");
    setProfileImageUrl("");
    setLivenessConfidence(0);
    sessionIdRef.current = "";
  };

  // ── Done screen ────────────────────────────────────────────────────────────

  if (phase === "done") {
    return (
      <div className="w-full rounded-2xl overflow-hidden shadow-[0_24px_80px_rgb(var(--brand-cold-rgb)/0.22)] border border-border/70 bg-card/85 backdrop-blur-sm">
        <div className="p-6">
          <div className="text-center">
            <div className="w-20 h-20 mx-auto mb-4 rounded-full bg-emerald-500/20 flex items-center justify-center">
              <CheckCircle2 className="w-10 h-10 text-[rgb(var(--brand-jungle-rgb))]" />
            </div>
            <h3 className="text-xl font-bold text-[rgb(var(--brand-jungle-rgb))] mb-2">Verification Successful!</h3>
            <p className="text-slate-400 text-sm mb-6">Your identity has been verified</p>

            {/* Score grid */}
            <div className="grid grid-cols-2 gap-4 text-left mb-6">
              <div className="bg-secondary/50 rounded-lg p-3">
                <p className="text-xs text-slate-400 mb-1">Liveness Score</p>
                <p className="text-2xl font-bold text-white">
                  {livenessConfidence > 0
                    ? `${(livenessConfidence * 100).toFixed(1)}%`
                    : livenessScore > 0
                    ? `${(livenessScore * 100).toFixed(1)}%`
                    : "—"}
                </p>
              </div>
              <div className="bg-secondary/50 rounded-lg p-3">
                <p className="text-xs text-slate-400 mb-1">Profile Image</p>
                <p className="text-xs text-[rgb(var(--brand-jungle-rgb))] truncate">
                  {profileImageUrl ? "Uploaded ✓" : "Pending"}
                </p>
              </div>
            </div>

            {/* Captured image */}
            {capturedImageB64 && (
              <div className="text-left mb-4">
                <p className="text-xs text-slate-400 mb-2">Captured Image</p>
                {/* eslint-disable-next-line @next/next/no-img-element */}
                <img
                  src={`data:image/jpeg;base64,${capturedImageB64}`}
                  alt="Captured"
                  className="w-full rounded-lg border border-slate-600/50 object-cover"
                />
              </div>
            )}

            {/* CDN URL */}
            {profileImageUrl && (
              <div className="text-left mb-4">
                <p className="text-xs text-slate-400 mb-1">CDN URL</p>
                <a
                  href={profileImageUrl}
                  target="_blank"
                  rel="noreferrer"
                  className="text-xs text-[rgb(var(--brand-jungle-rgb))] break-all underline"
                >
                  {profileImageUrl}
                </a>
              </div>
            )}
          </div>

          <button
            onClick={resetToIdle}
            className="auth-submit mt-2 w-full py-3 rounded-xl font-medium text-white transition-all"
          >
            Verify Again
          </button>
        </div>
      </div>
    );
  }

  // ── Error screen ───────────────────────────────────────────────────────────

  if (phase === "error") {
    return (
      <div className="w-full rounded-2xl overflow-hidden shadow-[0_24px_80px_rgb(var(--brand-cold-rgb)/0.22)] border border-border/70 bg-card/85 backdrop-blur-sm">
        <div className="p-6 text-center">
          <div className="w-20 h-20 mx-auto mb-4 rounded-full bg-red-500/20 flex items-center justify-center">
            <XCircle className="w-10 h-10 text-red-400" />
          </div>
          <h3 className="text-xl font-bold text-red-400 mb-2">Verification Failed</h3>
          <p className="text-slate-400 text-sm mb-6">{errorMsg || "Please try again"}</p>
          <button
            onClick={resetToIdle}
            className="auth-submit w-full py-3 rounded-xl font-medium text-white transition-all"
          >
            Try Again
          </button>
        </div>
      </div>
    );
  }

  // ── Main UI ────────────────────────────────────────────────────────────────

  // Brand gradient — matches portal's brand-btn-gradient
  const brandGradient = "[background:linear-gradient(120deg,rgb(var(--brand-jungle-rgb)/1)_0%,rgb(var(--brand-cold-rgb)/0.96)_100%)]";

  return (
    <div className="w-full rounded-2xl overflow-hidden shadow-[0_24px_80px_rgb(var(--brand-cold-rgb)/0.22)] border border-border/70 bg-card/85 backdrop-blur-sm">

      {/* ── Video section ──────────────────────────────────────────────── */}
      <div className="relative aspect-video bg-black">
        <video ref={videoRef} className="w-full h-full object-cover" autoPlay muted playsInline />
        {/* Overlay canvas — drawn from processFrame, NOT from RAF */}
        <canvas ref={overlayRef} className="absolute inset-0 w-full h-full pointer-events-none" />
        {/* Hidden capture canvas */}
        <canvas ref={canvasRef} className="hidden" />

        {/* Idle / starting overlay */}
        {(phase === "idle" || phase === "starting") && (
          <div className="absolute inset-0 flex items-center justify-center bg-black/60 backdrop-blur-sm">
            <div className="text-center px-6">
              <div className="w-24 h-24 mx-auto mb-4 rounded-full [background:linear-gradient(120deg,rgb(var(--brand-jungle-rgb)/1)_0%,rgb(var(--brand-cold-rgb)/0.96)_100%)] flex items-center justify-center shadow-lg">
                <Camera className="w-12 h-12 text-white" />
              </div>
              <h3 className="text-xl font-bold text-white mb-2">eKYC Verification</h3>
              <p className="text-slate-300 text-sm mb-6">Complete face liveness challenges to verify your identity</p>
              <button
                onClick={startEKYC}
                disabled={phase === "starting"}
                className="auth-submit px-8 py-3 rounded-xl font-semibold text-white transition-all shadow-lg disabled:opacity-50 flex items-center gap-2 mx-auto"
              >
                {phase === "starting" ? (
                  <><Loader2 className="w-5 h-5 animate-spin" />Starting…</>
                ) : (
                  <><Camera className="w-5 h-5" />Begin Verification</>
                )}
              </button>
            </div>
          </div>
        )}

        {/* Completing overlay */}
        {phase === "completing" && (
          <div className="absolute inset-0 flex flex-col items-center justify-center bg-black/75 backdrop-blur-sm gap-3">
            <Loader2 className="w-10 h-10 text-[rgb(var(--brand-jungle-rgb))] animate-spin" />
            <p className="text-white font-medium">Finalising verification…</p>
          </div>
        )}

        {/* Challenge instruction pill (active) */}
        {phase === "active" && challengeCfg && (
          <div className="absolute top-4 left-1/2 -translate-x-1/2 bg-black/70 backdrop-blur-sm px-6 py-3 rounded-full flex items-center gap-3 shadow-md whitespace-nowrap">
            <span className="text-2xl">{challengeCfg.icon}</span>
            <span className="text-white text-sm font-medium">
              {guidance[0] || challengeCfg.description}
            </span>
          </div>
        )}

        {/* Stats — top right */}
        {phase === "active" && (
          <div className="absolute top-4 right-4 bg-black/50 backdrop-blur-sm px-3 py-1.5 rounded-full text-xs text-white font-mono">
            {Math.round(1000 / FRAME_INTERVAL_MS)} FPS · {Math.round(avgLatency)}ms
          </div>
        )}

        {/* Face + liveness badges — bottom left */}
        {phase === "active" && (
          <div className="absolute bottom-8 left-4 flex gap-2 flex-wrap">
            <span className={`text-xs px-3 py-1 rounded-full font-medium ${
              faceDetected ? "bg-[rgb(var(--brand-jungle-rgb)/0.85)] text-white" : "bg-yellow-500/80 text-black"
            }`}>
              {faceDetected ? "Face detected" : "No face"}
            </span>
            {livenessScore > 0 && (
              <span className="text-xs px-3 py-1 rounded-full bg-black/60 text-white">
                Liveness {Math.round(livenessScore * 100)}%
              </span>
            )}
          </div>
        )}

        {/* Progress bar — bottom edge */}
        {phase === "active" && (
          <div className="absolute bottom-0 left-0 right-0 h-1 bg-border/40">
            <div
              className="h-full [background:linear-gradient(120deg,rgb(var(--brand-jungle-rgb)/1)_0%,rgb(var(--brand-cold-rgb)/0.96)_100%)] transition-all duration-300"
              style={{ width: `${Math.round(overallProgress * 100)}%` }}
            />
          </div>
        )}
      </div>

      {/* ── Step cards — always visible (idle shows DEFAULT_STEPS as grey) ── */}
      <div className="p-4 bg-secondary/30 border-t border-border/60">
        <div className="flex justify-between gap-2">
          {(steps.length > 0 ? steps : DEFAULT_STEPS).map((s, i) => {
            // Handle both full enum ("EKYC_STEP_COMPLETED") and stripped ("COMPLETED")
            const rawState = s.state ?? "";
            const isDone   = rawState === "EKYC_STEP_COMPLETED"  || rawState === "COMPLETED";
            const isActive = rawState === "EKYC_STEP_IN_PROGRESS" || rawState === "IN_PROGRESS";
            const isFail   = rawState === "EKYC_STEP_FAILED"     || rawState === "FAILED";

            // Type: try full enum first, then stripped
            const rawType = s.type ?? s.challenge_type ?? "";
            const cfg = CHALLENGE_CONFIG[rawType] || CHALLENGE_CONFIG[stripPrefix(rawType)] || { icon: "🔍", label: rawType || "Step", description: "" };

            return (
              <div
                key={i}
                className={`flex-1 p-3 rounded-lg border transition-all ${
                  isDone   ? "bg-[var(--brand-surface-3)] border-[rgb(var(--brand-jungle-rgb)/0.4)]" :
                  isActive ? "bg-[var(--brand-surface-2)] border-[rgb(var(--brand-cold-rgb)/0.4)] animate-pulse" :
                  isFail   ? "bg-destructive/10 border-destructive/30" :
                             "bg-secondary/50 border-border/50"
                }`}
              >
                <div className="flex items-center gap-2 mb-1">
                  <span className="text-lg">{cfg.icon}</span>
                  {isDone && <CheckCircle2 className="w-4 h-4 text-[rgb(var(--brand-jungle-rgb))]" />}
                  {isFail && <XCircle      className="w-4 h-4 text-destructive" />}
                </div>
                <p className="text-xs text-foreground/80 font-medium">{cfg.label}</p>
              </div>
            );
          })}
        </div>
      </div>

      {/* ── Inline error strip ─────────────────────────────────────────── */}
      {inlineError && phase === "active" && (
        <div className="px-4 py-3 bg-destructive/10 border-t border-destructive/20 flex items-start gap-3">
          <AlertCircle className="w-5 h-5 text-destructive shrink-0 mt-0.5" />
          <div>
            <p className="text-sm font-medium text-destructive">Detection issue</p>
            <p className="text-xs text-destructive/80">{inlineError}</p>
          </div>
        </div>
      )}

      {/* ── Cancel button (active) ─────────────────────────────────────── */}
      {phase === "active" && (
        <div className="p-4 bg-secondary/20 border-t border-border/60">
          <button
            onClick={() => { stopCamera(); resetToIdle(); }}
            className="w-full py-2.5 bg-destructive/10 hover:bg-destructive/20 border border-destructive/20 rounded-xl text-sm font-medium text-destructive transition-all flex items-center justify-center gap-2"
          >
            <CameraOff className="w-5 h-5" />
            Cancel Verification
          </button>
        </div>
      )}

      {/* ── How it works (idle only) ───────────────────────────────────── */}
      {phase === "idle" && (
        <div className="p-4 bg-secondary/20 border-t border-border/60">
          <p className="text-xs font-semibold text-[rgb(var(--brand-cold-rgb))] mb-2 uppercase tracking-wide">How it works</p>
          <div className="space-y-1.5 text-xs text-muted-foreground">
            <p>🌟 Ensure you are in a well-lit area</p>
            <p>👤 Position your face in the centre of the oval</p>
            <p>👓 Remove glasses if possible for better detection</p>
            <p>👁️ <strong className="text-foreground">Blink</strong> → <strong className="text-foreground">Look Left</strong> → <strong className="text-foreground">Look Right</strong> → <strong className="text-foreground">Look Ahead</strong></p>
            <p>⏱️ The process takes about 30 seconds</p>
          </div>
        </div>
      )}
    </div>
  );
}
