import type { Metadata } from "next";
import { GeistSans } from "geist/font/sans";
import { GeistMono } from "geist/font/mono";
import "./globals.css";

import { getServerSession } from "@lib/auth/session";
import { PortalSessionProvider } from "@lib/auth/portal-session-context";

export const metadata: Metadata = {
  title: "Labaid Insuretech B2B Dashboard",
  description:
    "Labaid Insuretech B2B Dashboard built with Next.js and Tailwind CSS",
  icons: {
    icon: [
      {
        url: "logos/favicon.svg",
        media: "(prefers-color-scheme: light)",
      },
      {
        url: "logos/favicon-dark.svg",
        media: "(prefers-color-scheme: dark)",
      },
      {
        url: "logos/favicon.svg",
        type: "image/svg+xml",
      },
    ],
    apple: "logos/favicon.svg",
  },
};

export default async function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const session = await getServerSession();

  return (
    <html lang="en">
      <body
        className={`${GeistSans.variable} ${GeistMono.variable} antialiased`}
      >
        <PortalSessionProvider initialSession={session}>
          {children}
        </PortalSessionProvider>
      </body>
    </html>
  );
}
