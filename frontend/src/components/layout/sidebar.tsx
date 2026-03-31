"use client";

import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { cn } from "@/lib/utils";
import { clearToken, getUser } from "@/lib/auth";
import { Button } from "@/components/ui/button";

const navItems = [
  { href: "/sources", label: "货源管理", icon: "📦", roles: null as string[] | null },
  { href: "/products", label: "商品管理", icon: "🏷️", roles: null },
  { href: "/publish", label: "发品任务", icon: "🚀", roles: null },
  { href: "/users", label: "用户管理", icon: "👥", roles: ["super_admin", "admin"] },
];

export function Sidebar() {
  const pathname = usePathname();
  const router = useRouter();
  const user = getUser();
  const userRoles: string[] = user?.roles || [];

  const handleLogout = () => {
    clearToken();
    router.replace("/login");
  };

  const visibleItems = navItems.filter(
    (item) => !item.roles || item.roles.some((r) => userRoles.includes(r))
  );

  return (
    <aside className="w-56 border-r bg-gray-50 p-4 min-h-screen flex flex-col">
      <h1 className="text-lg font-bold mb-6 px-2">代发工具</h1>
      <nav className="space-y-1 flex-1">
        {visibleItems.map((item) => (
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
      {user && (
        <div className="border-t pt-4 mt-4">
          <div className="px-2 text-sm text-gray-600 mb-2 truncate">{user.email}</div>
          <Button variant="outline" size="sm" className="w-full" onClick={handleLogout}>
            退出登录
          </Button>
        </div>
      )}
    </aside>
  );
}
