"use client";

import { usePathname, useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { isLoggedIn } from "@/lib/auth";
import { Sidebar } from "./sidebar";

const publicPaths = ["/login", "/login/callback"];

export function AuthGuard({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const router = useRouter();
  const [checked, setChecked] = useState(false);
  const isPublic = publicPaths.some((p) => pathname.startsWith(p));

  useEffect(() => {
    if (!isPublic && !isLoggedIn()) {
      router.replace("/login");
    } else {
      setChecked(true);
    }
  }, [pathname, isPublic, router]);

  if (!checked) return null;

  if (isPublic) {
    return <>{children}</>;
  }

  return (
    <div className="flex">
      <Sidebar />
      <main className="flex-1 p-6">{children}</main>
    </div>
  );
}
