import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import { AuthGuard } from "@/components/layout/auth-guard";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "代发工具",
  description: "电商一键代发运营工具",
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="zh-CN">
      <body className={inter.className}>
        <AuthGuard>{children}</AuthGuard>
      </body>
    </html>
  );
}
