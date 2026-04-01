# eKYC UI Components Report - B2B Portal

## Executive Summary
The B2B portal implements a complete eKYC (Electronic Know Your Customer) flow with:
- **Live video capture** with canvas-based overlay rendering
- **Real-time face detection & liveness scoring** with visual feedback
- **Multi-step challenge flow** (Blink, Look Left, Look Right, Capture)
- **Dynamic progress tracking** and guidance messages
- **Head pose indicator** showing yaw/pitch/roll in real-time

---

## File Inventory

### Core eKYC Components

| File | Location | Purpose |
|------|----------|---------|
| `EKYCFlow.tsx` | `components/kyc/` | Main eKYC UI component with camera, canvas overlay, step tracking |
| `kyc-page-client.tsx` | `app/kyc/` | Client wrapper that manages KYC cookie state |
| `page.tsx` | `app/kyc/` | Server-side KYC gate (checks if already verified) |

### Backend API Routes

| Route | File | Purpose |
|-------|------|---------|
| `POST /api/auth/kyc/initiate` | `app/api/auth/kyc/initiate/route.ts` | Creates KYC session, returns steps & session_id |
| `POST /api/auth/kyc/frame` | `app/api/auth/kyc/frame/route.ts` | Forwards video frames to gateway for processing |
| `POST /api/auth/kyc/complete` | `app/api/auth/kyc/complete/route.ts` | Finalizes KYC, sets status to PENDING_REVIEW |

---

## Component Architecture: EKYCFlow.tsx

### State Management

**React State (UI-driven):**
- `phase`: "idle" | "starting" | "active" | "completing" | "done" | "error"
- `sessionId`: Current eKYC session ID
- `steps`: Array of KYC challenges (Blink, Look Left, etc.)
- `currentStepIdx`: Index of current step
- `overallProgress`: 0-1 progress value
- `guidance`: Array of guidance messages
- `faceDetected`: Boolean face detection status
- `livenessScore`: 0-1 liveness confidence
- `frameCount`: Number of frames processed
- `avgLatency`: Average round-trip latency (ms)

**Ref-based State (frame loop control):**
- `sessionIdRef`: Session ID (avoids stale closure in frame loop)
- `isRunningRef`: Boolean flag to control frame loop
- `sessionStateRef`: Last session state from gateway
- `frameDataRef`: Latest frame response data

### Canvas Overlay Drawing System

**Face Guide Oval** (Always Visible)
```typescript
// Centered, 40% width × 70% height
const centerX = overlay.width / 2;
const centerY = overlay.height / 2;
const ovalWidth = overlay.width * 0.4;
const ovalHeight = overlay.height * 0.7;

ctx.strokeStyle = detected ? "rgba(34,197,94,0.8)" : "rgba(255,255,255,0.3)";
ctx.lineWidth = 2;
ctx.setLineDash([12, 6]); // Dashed line pattern
ctx.beginPath();
ctx.ellipse(centerX, centerY, ovalWidth / 2, ovalHeight / 2, 0, 0, 2 * Math.PI);
ctx.stroke();
```
- **Green when face detected**, white when not detected
- Acts as frame guide for user positioning

**Bounding Box** (Detection-based)
```typescript
const box = data.detection?.box; // Normalized 0-1 coords
if (detected && box) {
  const liveness = data.liveness_score ?? data.liveness_confidence ?? 0;
  ctx.strokeStyle = liveness > 0.5 ? "#22c55e" : "#ef4444"; // Green or red
  ctx.lineWidth = 2;
  ctx.strokeRect(
    box.x * scaleX,
    box.y * scaleY,
    box.width * scaleX,
    box.height * scaleY
  );
}
```
- **Green box**: Liveness score > 0.5 (high confidence)
- **Red box**: Low liveness confidence
- Coordinates are normalized from gateway, scaled to canvas

**Eye Contours** (Detected when blinking/opening)
```typescript
const drawEye = (eye?: EyeContour) => {
  if (!eye?.edges || !eye?.points) return;
  const eyeColor = isBlinking ? "rgba(34,197,94,0.9)" : "rgba(59,130,246,0.85)";
  ctx.strokeStyle = eyeColor;
  ctx.lineWidth = 1.5;
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
```
- **Green contour**: When blinking detected
- **Blue contour**: When eyes open
- Contour data includes edge connections and landmark points
- Coordinates scaled from video dimensions to overlay

**Head Pose Indicator** (Top-right corner)
```typescript
const hp = data.head_pose; // { yaw, pitch, roll }
const ix = overlay.width - 70;
const iy = 10;

// Dark background
ctx.fillStyle = "rgba(0,0,0,0.5)";
ctx.fillRect(ix - 30, iy, 60, 60);

// White circle outline
ctx.strokeStyle = "#ffffff";
ctx.lineWidth = 2;
ctx.beginPath();
ctx.arc(ix, iy + 30, 20, 0, 2 * Math.PI);
ctx.stroke();

// Green dot (position = center + yaw/pitch offset)
const yawRad = (hp.yaw * Math.PI) / 180;
const pitchRad = (hp.pitch * Math.PI) / 180;
const dx = Math.sin(yawRad) * 20;
const dy = Math.sin(pitchRad) * 20;

ctx.fillStyle = "#22c55e";
ctx.beginPath();
ctx.arc(ix + dx, iy + 30 + dy, 5, 0, 2 * Math.PI);
ctx.fill();

// Label
ctx.fillStyle = "#ffffff";
ctx.font = "10px sans-serif";
ctx.textAlign = "center";
ctx.fillText("Pose", ix, iy + 55);
```
- Compass-like indicator: white circle with green dot showing head direction
- Green dot moves based on yaw (left-right) and pitch (up-down)
- Helps user correct head position in real-time

### Canvas Rendering Loop

**Overlay Drawing** (Runs on requestAnimationFrame)
```typescript
useEffect(() => {
  if (phase !== "active") return;
  let rafId: number;
  const loop = () => {
    drawOverlay();
    rafId = requestAnimationFrame(loop);
  };
  rafId = requestAnimationFrame(loop);
  return () => cancelAnimationFrame(rafId);
}, [phase, drawOverlay]);
```
- Continuous smooth 60fps rendering of overlay
- Updates dynamically as frameDataRef changes
- Stops when phase leaves "active"

### Frame Processing Pipeline

**Timing & Latency Tracking:**
```typescript
const t0 = performance.now();
const resp = await fetch("/api/auth/kyc/frame", { /* ... */ });
const latency = performance.now() - t0;
latenciesRef.current.push(latency);
if (latenciesRef.current.length > 30) latenciesRef.current.shift();
setAvgLatency(/* average of last 30 */)
```
- Maintains rolling window of 30 frame latencies
- Displays avg latency in top-right (e.g., "42 frames · 234ms")

**Frame Interval:**
- `FRAME_INTERVAL_MS = 100` → 10 FPS capture rate (100ms between submissions)
- Uses `setTimeout` with `isRunningRef` flag (not React state) to avoid stale closures
- Only reschedules if session state is still `EKYC_SESSION_ACTIVE`

---

## UI Layout & Tailwind Classes

### Camera Card Container
```tsx
<div className="rounded-2xl overflow-hidden bg-slate-900 shadow-2xl">
  <div className="relative aspect-video bg-black">
    {/* Video + Overlay Canvas + Hidden Processing Canvas */}
  </div>
</div>
```
- **`bg-slate-900`**: Dark background matching eKYC ambiance
- **`shadow-2xl`**: Heavy drop shadow for depth
- **`aspect-video`**: 16:9 ratio for video

### Video & Overlay Layers
```tsx
<video ref={videoRef} className="w-full h-full object-cover" />
<canvas ref={overlayRef} className="absolute inset-0 w-full h-full pointer-events-none" />
<canvas ref={canvasRef} className="hidden" /> {/* Frame capture, not displayed */}
```
- Video: fills container, covers object
- Overlay: absolute positioned over video, pointer-events-none
- Processing canvas: hidden, used for JPEG compression

### Idle/Starting Overlay
```tsx
{(phase === "idle" || phase === "starting") && (
  <div className="absolute inset-0 flex flex-col items-center justify-center bg-black/70 backdrop-blur-sm gap-4">
    <div className="w-24 h-24 rounded-full bg-gradient-to-br from-emerald-500 to-cyan-500 ...">
      <span className="text-5xl">📷</span>
    </div>
    <h3 className="text-xl font-bold text-white">eKYC Verification</h3>
    <Button className="bg-gradient-to-r from-emerald-500 to-cyan-500 ...">
      Begin Verification
    </Button>
  </div>
)}
```
- **Color Scheme**: Emerald (green) to Cyan (blue) gradient
- **Backdrop blur**: Semi-transparent with blur effect
- **Gradient Button**: Matches brand color scheme

### Active Phase UI Elements

**Challenge Instruction** (Top center)
```tsx
{phase === "active" && (
  <div className="absolute top-4 left-1/2 -translate-x-1/2 bg-black/70 backdrop-blur-sm px-6 py-3 rounded-full flex items-center gap-3">
    <span className="text-2xl">{challengeCfg.icon}</span>
    <span className="text-white text-sm font-medium">{guidance[0]}</span>
  </div>
)}
```
- Displays current challenge icon (👁️ Blink, 👈 Look Left, etc.)
- Shows real-time guidance message from gateway

**FPS & Latency Display** (Top right)
```tsx
{phase === "active" && (
  <div className="absolute top-4 right-4 bg-black/50 backdrop-blur-sm px-3 py-1.5 rounded-full text-xs text-white">
    {frameCount} frames · {Math.round(avgLatency)}ms
  </div>
)}
```

**Face Detection & Liveness Status** (Bottom left)
```tsx
{phase === "active" && (
  <div className="absolute bottom-4 left-4 flex gap-2 flex-wrap">
    <span className={faceDetected ? "bg-emerald-500/80 text-white" : "bg-yellow-500/80 text-black"}>
      {faceDetected ? "Face detected" : "No face"}
    </span>
    {livenessScore > 0 && (
      <span className="text-xs px-3 py-1 rounded-full bg-black/60 text-white">
        Liveness: {Math.round(livenessScore * 100)}%
      </span>
    )}
  </div>
)}
```

**Progress Bar** (Bottom, full width)
```tsx
{phase === "active" && (
  <div className="absolute bottom-0 left-0 right-0 h-1 bg-slate-700/80">
    <div
      className="h-full bg-gradient-to-r from-emerald-500 to-cyan-500 transition-all duration-300"
      style={{ width: `${Math.round(overallProgress * 100)}%` }}
    />
  </div>
)}
```
- Thin 1px bar at bottom
- Gradient fill from emerald to cyan
- Smooth 300ms transition animation

### Step Tracker Bar
```tsx
{steps.length > 0 && (
  <div className="p-4 bg-slate-800/60 flex justify-between gap-2">
    {steps.map((s, i) => (
      <div className={`flex-1 p-3 rounded-xl border transition-all text-center ${
        isCompleted ? "bg-emerald-500/20 border-emerald-500/50 text-emerald-300"
        : isInProgress ? "bg-cyan-500/20 border-cyan-500/50 text-cyan-200 animate-pulse"
        : isFailed ? "bg-red-500/20 border-red-500/50 text-red-300"
        : "bg-slate-700/40 border-slate-600/40 text-slate-400"
      }`}>
        <div className="text-xl mb-1">{cfg.icon}</div>
        <div className="text-xs font-medium">{cfg.label}</div>
        <div className="text-xs mt-1">{isCompleted ? "✓" : isInProgress ? "…" : isFailed ? "✗" : "○"}</div>
      </div>
    ))}
  </div>
)}
```

**States:**
- **Completed** (✓): `bg-emerald-500/20`, emerald border & text
- **In Progress** (…): `bg-cyan-500/20`, cyan border & text, **animated pulse**
- **Failed** (✗): `bg-red-500/20`, red border & text
- **Pending** (○): `bg-slate-700/40`, gray border & text

### Guidance Card
```tsx
{phase === "active" && guidance.length > 0 && (
  <div className="bg-blue-500/10 border border-blue-500/30 rounded-xl p-3">
    <div className="flex items-start gap-2">
      <span className="text-blue-400 text-sm">💬</span>
      {guidance.map(g => <p className="text-sm text-blue-300">{g}</p>)}
    </div>
  </div>
)}
```
- Light blue background with semi-transparent border
- Lists all current guidance messages from backend

### Success & Error Screens

**Done Phase:**
```tsx
<Card className="max-w-md mx-auto text-center border-emerald-200">
  <div className="w-20 h-20 mx-auto mb-4 rounded-full bg-emerald-500/20 flex items-center justify-center">
    <span className="text-5xl">✅</span>
  </div>
  <CardTitle className="text-emerald-600">Verification Successful</CardTitle>
  <img src={profileImageUrl} className="rounded-full w-24 h-24 mx-auto object-cover border-4 border-emerald-200" />
</Card>
```

**Error Phase:**
```tsx
<Card className="max-w-md mx-auto text-center border-red-200">
  <div className="w-20 h-20 mx-auto mb-4 rounded-full bg-red-500/20 flex items-center justify-center">
    <span className="text-5xl">❌</span>
  </div>
  <CardTitle className="text-red-600">Verification Failed</CardTitle>
  <Button onClick={() => { /* reset and retry */ }}>Try Again</Button>
</Card>
```

---

## Global Tailwind & CSS Configuration

### Color Palette (from `globals.css`)

**Brand Colors:**
```css
--brand-jungle: #03a765;      /* Primary green */
--brand-cold: #123f50;        /* Secondary blue */
--brand-dark-grey: #a8a9ad;   /* Neutral */
```

**Semantic Colors:**
```css
--primary: #03a765;           /* Brand green */
--accent: #123f50;            /* Brand blue */
--success: #03a765;           /* Green */
--warning: #a8a9ad;           /* Gray */
--destructive: oklch(0.577 0.245 27.325); /* Red */
--alert-critical: #d9534f;
--alert-info: #03a765;
```

**eKYC-specific colors used:**
- **Emerald-500** (`#10b981`): Face detected, step completed
- **Cyan-500** (`#06b6d4`): In-progress challenges
- **Yellow-500** (`#eab308`): No face detected
- **Red-500** (`#ef4444`): Low liveness, failed steps
- **Blue-500** (`#3b82f6`): Guidance messages, eye contours

### Component Classes (from `globals.css`)

```css
.portal-shell { @apply min-h-screen bg-background; }
.portal-panel { @apply overflow-hidden rounded-xl border border-border/90 bg-card; }
.auth-submit {
  @apply h-11 w-full rounded-lg text-sm font-semibold text-white transition-all;
  background: linear-gradient(120deg, rgb(var(--brand-jungle-rgb) / 1) 0%, rgb(var(--brand-cold-rgb) / 0.95) 100%);
}
```

---

## API Integration Flow

### 1. Initiate Session
**Endpoint:** `POST /api/auth/kyc/initiate`
**Request:** `{ user_id: string }`
**Response:**
```json
{
  "ok": true,
  "kyc_verification_id": "uuid",
  "session_id": "FLVE-session-id",
  "steps": [
    { "step_number": 1, "type": "EKYC_CHALLENGE_BLINK", "state": "PENDING", "instruction": "Blink your eyes" },
    { "step_number": 2, "type": "EKYC_CHALLENGE_LOOK_LEFT", "state": "PENDING", "instruction": "Turn left" },
    ...
  ],
  "session_state": "EKYC_SESSION_ACTIVE"
}
```

### 2. Submit Frame (100ms interval)
**Endpoint:** `POST /api/auth/kyc/frame`
**Request:** multipart/form-data
```
user_id: string
session_id: string
file: Blob (JPEG, 640×480)
```
**Response:**
```json
{
  "session_state": "EKYC_SESSION_ACTIVE",
  "detection": {
    "detected": true,
    "box": { "x": 0.2, "y": 0.15, "width": 0.6, "height": 0.7 }
  },
  "head_pose": { "yaw": 5.2, "pitch": -2.1, "roll": 0.3 },
  "eye_state": { "left_openness": 0.95, "right_openness": 0.92, "is_blinking": false },
  "eye_contours": {
    "left": {
      "edges": [[0,1], [1,2], [2,3]],
      "points": {"0": {"x": 0.35, "y": 0.4}, "1": {...}, ...}
    },
    "right": { ... }
  },
  "liveness_score": 0.87,
  "liveness_confidence": 0.87,
  "overall_progress": 0.45,
  "guidance": ["Blink your eyes", "Look more towards the camera"],
  "current_step_detail": {
    "step_number": 1,
    "type": "EKYC_CHALLENGE_BLINK",
    "state": "IN_PROGRESS",
    "confidence": 0.45
  }
}
```

### 3. Complete Session
**Endpoint:** `POST /api/auth/kyc/complete`
**Request:** `{ session_id: string }`
**Response:**
```json
{
  "ok": true,
  "profile_image_url": "https://..."
}
```
**Side Effect:** Sets `portal_kyc_verified=pending_review` cookie (12-hour expiry)

---

## Canvas Drawing Technical Details

### Coordinate Systems

**Video Frame Coordinates:**
- Origin: top-left of video stream
- Range: [0, videoWidth] × [0, videoHeight]
- Gateway returns box/contours in **normalized 0-1 coordinates**

**Overlay Canvas Coordinates:**
- Origin: top-left of overlay canvas
- Range: [0, overlay.width] × [0, overlay.height]
- Scaling: `scaleX = overlay.width / videoWidth`, `scaleY = overlay.height / videoHeight`

**Transformation:**
```typescript
const scaleX = overlay.width / videoWidth;
const scaleY = overlay.height / videoHeight;

// Gateway normalized coords → overlay canvas coords
ctx.strokeRect(
  box.x * scaleX,
  box.y * scaleY,
  box.width * scaleX,
  box.height * scaleY
);
```

### Drawing Order (Z-order from back to front)
1. **Video layer** (bottom)
2. **Overlay canvas**:
   - Face guide oval (always)
   - Bounding box (when detected)
   - Eye contours (when available)
   - Head pose indicator (when available)
   - Guidance text overlays (text layer)
3. **UI elements** (text badges, progress bar)

### Rendering Performance
- **Frame capture**: 640×480 JPEG @ 85% quality, 100ms interval (10 FPS)
- **Overlay redraw**: 60 FPS (requestAnimationFrame)
- **Latency tracking**: Rolling average of last 30 frames
- **Memory**: Single video stream + overlay canvas + hidden processing canvas

---

## Face Contour & Eye Overlay Details

### Why Eye Contours Are Conditional
Eye contours are **only drawn when the gateway provides them** in `frameData.eye_contours`:
```typescript
if (detected && eyeContours) {
  const isBlinking = data.eye_state?.is_blinking ?? false;
  const eyeColor = isBlinking ? "rgba(34,197,94,0.9)" : "rgba(59,130,246,0.85)";
  drawEye(eyeContours.left);
  drawEye(eyeContours.right);
}
```

**Data Structure:**
```typescript
interface EyeContour {
  edges: [number, number][]; // Array of [pointA_idx, pointB_idx] pairs
  points: Record<string, { x: number; y: number }>; // Normalized coords
}
```

**Drawing Algorithm:**
1. For each edge pair (a, b) in `eye.edges`:
2. Lookup points[String(a)] and points[String(b)]
3. Draw line segment: `ctx.moveTo(pa.x * scaleX, pa.y * scaleY)` → `ctx.lineTo(pb.x * scaleX, pb.y * scaleY)`
4. Stroke with color based on blink state

### Why No Static Contour Drawing
- **No hardcoded contours**: The component does NOT draw static eye shapes
- **Gateway-driven**: Only draws what the backend provides
- **Dynamic updates**: Contours change every frame based on facial feature detection

---

## Progress Indication Mechanisms

### 1. Overall Progress Bar (Bottom)
- **Source**: `data.overall_progress` (0-1 float from gateway)
- **Display**: Width percentage, 300ms smooth transition
- **Color**: Emerald-to-cyan gradient
- **Height**: 1px (thin bar at bottom of video)

### 2. Step Tracker (Below video)
- **Source**: `steps` array with `state` field per step
- **States**:
  - `PENDING`: Gray circle (○)
  - `IN_PROGRESS`: Cyan with animated pulse (…)
  - `COMPLETED`: Green with checkmark (✓)
  - `FAILED`: Red with X (✗)
- **Visual**: 5-6 step boxes, each showing icon + state indicator

### 3. Liveness Score Badge (Bottom-left)
- **Source**: `data.liveness_score` or `data.liveness_confidence` (0-1)
- **Display**: Only shown if > 0
- **Format**: "Liveness: X%" (e.g., "Liveness: 87%")
- **Color**: Black/60 background with white text

### 4. Real-time Guidance Messages (Above video + separate card)
- **Source**: `data.guidance` or `data.guidance_messages` (string[])
- **Display Locations**:
  - **Top center**: First message only (in rounded pill)
  - **Guidance card**: All messages in blue-tinted box
- **Updates**: Every frame (100ms interval)

---

## Missing Features / Observations

### What's NOT Implemented
1. **Static mouth/nose/face outline drawing**: Only dynamic contours from backend
2. **Head tilt indicator** (separate from pose indicator): Only pose circle
3. **Liveness replay video**: No video recording or playback
4. **Document capture overlay**: No ID/passport detection
5. **Iris/pupil contours**: Only eye bounding contours (edges + points)
6. **Real-time audio/speaking detection**: No audio processing

### Deliberate Design Choices
1. **No video mirroring on display**: Uses raw camera coordinates (matches FLVE detection)
2. **Canvas overlay separate from video**: Allows independent redraw rates (60fps overlay, 10fps frames)
3. **Ref-based session control**: Avoids stale closure issues in frame loop
4. **Client-side KYC initiation**: Server-side initiate caused self-referential fetch loop in dev

---

## Summary Table: All Canvas Calls

| Element | Method | Color | Width | Pattern |
|---------|--------|-------|-------|---------|
| Face Guide Oval | `ellipse()` | Green if detected, white if not | 2px | Dashed [12,6] |
| Bounding Box | `strokeRect()` | Green (liveness>0.5) or red | 2px | Solid |
| Eye Left Edge | `lineTo()` (multiple) | Green (blinking) or blue | 1.5px | Solid |
| Eye Right Edge | `lineTo()` (multiple) | Green (blinking) or blue | 1.5px | Solid |
| Pose Circle | `arc()` | White outline | 2px | Solid |
| Pose Indicator Dot | `arc()` filled | Green (#22c55e) | 5px radius | Solid |
| Pose Background | `fillRect()` | Black/50% opacity | — | Solid |
| Pose Label Text | `fillText()` | White | — | 10px sans-serif |

---

## Tailwind Classes Used in EKYCFlow

**Container Classes:**
- `max-w-2xl`, `mx-auto`, `space-y-4`
- `rounded-2xl`, `overflow-hidden`, `shadow-2xl`
- `relative`, `aspect-video`, `bg-black`

**Color & Background:**
- `bg-slate-900`, `bg-black`, `bg-black/70`, `bg-black/50`
- `bg-emerald-500/20`, `bg-cyan-500/20`, `bg-red-500/20`
- `bg-yellow-500/80`, `text-white`, `text-black`

**Typography:**
- `text-xs`, `text-sm`, `text-xl`, `text-5xl`
- `font-bold`, `font-semibold`, `font-medium`

**Positioning & Layout:**
- `absolute`, `inset-0`, `left-1/2`, `-translate-x-1/2`
- `top-4`, `right-4`, `bottom-4`, `left-4`, `bottom-0`
- `flex`, `flex-col`, `items-center`, `justify-center`
- `gap-2`, `gap-3`, `gap-4`

**Effects:**
- `backdrop-blur-sm`, `pointer-events-none`
- `rounded-full`, `rounded-xl`
- `animate-pulse`
- `transition-all`, `duration-300`

**Badges & Buttons:**
- `px-3`, `py-1.5`, `px-6`, `py-3`
- `border`, `border-emerald-200`, `border-red-200`
- `shadow-lg`, `shadow-2xl`

---

**Report Generated**: Full eKYC UI component analysis
**Last Updated**: Current B2B Portal Codebase
