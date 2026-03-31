"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { cn } from "@/lib/utils";

const navItems = [
  { href: "/sources", label: "货源管理", icon: "📦" },
  { href: "/products", label: "商品管理", icon: "🏷️" },
  { href: "/publish", label: "发品任务", icon: "🚀" },
];

export function Sidebar() {
  const pathname = usePathname();

  return (
    <aside className="w-56 border-r bg-gray-50 p-4 min-h-screen">
      <h1 className="text-lg font-bold mb-6 px-2">代发工具</h1>
      <nav className="space-y-1">
        {navItems.map((item) => (
          <Link
            key={item.href}
            href={item.href}
            className={cn(
              "flex items-center gap-2 px-3 py-2 rounded-md text-sm transition-colors",
              pathname.startsWith(item.href)
                ? "bg-white shadow-sm font-medium"
                : "text-gray-600 hover:bg-white/60"
            )}
          >
            <span>{item.icon}</span>
            {item.label}
          </Link>
        ))}
      </nav>
    </aside>
  );
}
