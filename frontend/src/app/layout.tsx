import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";
import { ThemeProvider } from '@/components/theme/ThemeProvider';
import ErrorBoundary from '@/components/ErrorBoundary';

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: {
    default: "TweetVault - Your Twitter Bookmarks Archive",
    template: "%s | TweetVault"
  },
  description: "Organize and manage your Twitter bookmarks with TweetVault. Search, tag, and archive your favorite tweets.",
  keywords: ["Twitter", "Bookmarks", "Archive", "Social Media", "Organization"],
  authors: [{ name: "TweetVault" }],
  creator: "TweetVault",
  openGraph: {
    type: "website",
    locale: "en_US",
    url: "https://tweetvault.app",
    title: "TweetVault - Your Twitter Bookmarks Archive",
    description: "Organize and manage your Twitter bookmarks with TweetVault. Search, tag, and archive your favorite tweets.",
    siteName: "TweetVault",
  },
  twitter: {
    card: "summary_large_image",
    title: "TweetVault - Your Twitter Bookmarks Archive",
    description: "Organize and manage your Twitter bookmarks with TweetVault",
    creator: "@tweetvault",
  },
  icons: {
    icon: [
      {
        url: "/favicon.ico",
        sizes: "any",
      },
      {
        url: "/icon.svg",
        type: "image/svg+xml",
      },
    ],
    apple: [
      {
        url: "/apple-touch-icon.png",
        sizes: "180x180",
      },
    ],
  },
  manifest: "/site.webmanifest",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased`}
        suppressHydrationWarning
      >
        <ErrorBoundary>
          <ThemeProvider>{children}</ThemeProvider>
        </ErrorBoundary>
      </body>
    </html>
  );
}
