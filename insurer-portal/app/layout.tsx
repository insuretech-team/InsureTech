import type { Metadata } from "next";
import { IBM_Plex_Sans, Sora } from "next/font/google";

import "./globals.css";

const heading = Sora({
  subsets: ["latin"],
  variable: "--font-heading",
});

const body = IBM_Plex_Sans({
  subsets: ["latin"],
  variable: "--font-body",
  weight: ["400", "500", "600", "700"],
});

export const metadata: Metadata = {
  title: "LabaidInsuretech Insurer Dashboard",
  description: "LabaidInsuretech insurer dashboard for proposals, claims, products, and operational configuration.",
  icons: {
    icon: "/logos/favicon.svg",
    shortcut: "/logos/favicon.svg",
    apple: "/logos/favicon.svg",
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className={`${heading.variable} ${body.variable}`}>{children}</body>
    </html>
  );
}
