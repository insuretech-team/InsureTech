import type { Metadata } from "next";
import { GeistSans } from "geist/font/sans";
import { GeistMono } from "geist/font/mono";
import "./globals.css";

export const metadata: Metadata = {
  title: "Labaid Insuretech Partner Portal",
  description:
    "Labaid Insuretech Partner Portal built with Next.js and Tailwind CSS",
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

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body
        className={`${GeistSans.variable} ${GeistMono.variable} antialiased`}
      >
        {children}
      </body>
    </html>
  );
}
