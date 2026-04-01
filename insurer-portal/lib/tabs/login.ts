export const loginPageCopy = {
  header: {
    eyebrow: "Insurer access",
    title: "Sign in to the operations desk",
  },
  form: {
    mobileLabel: "Mobile number",
    mobilePlaceholder: "+8801712345678",
    passwordLabel: "Password",
    passwordPlaceholder: "Enter your password",
  },
  messages: {
    invalidCredentials: "Unable to sign in. Please verify your mobile number and password.",
    serviceDown: "The portal could not reach the authentication service.",
  },
  submit: {
    idle: "Sign in",
    busy: "Signing in...",
  },
  footer: "Sign in with your insurer account to continue to the dashboard.",
} as const;
