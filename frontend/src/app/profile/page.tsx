"use client";

import { useEffect, useState } from "react";
import { api } from "@/lib/api";
import { Badge } from "@/components/ui/badge";

export default function ProfilePage() {
  const [user, setUser] = useState<any>(null);

  useEffect(() => {
    api.auth.me().then(setUser).catch(() => {});
  }, []);

  if (!user) return <div className="p-8">加载中...</div>;

  return (
    <div className="max-w-lg">
      <h2 className="text-2xl font-bold mb-6">个人中心</h2>
      <div className="bg-white rounded-xl shadow-sm border p-6 space-y-4">
        <div className="flex items-center gap-4">
          {user.avatar ? (
            <img src={user.avatar} alt="avatar" className="w-16 h-16 rounded-full" />
          ) : (
            <div className="w-16 h-16 rounded-full bg-gray-200 flex items-center justify-center text-2xl">
              {(user.name || user.email || "U").charAt(0).toUpperCase()}
            </div>
          )}
          <div>
            <h3 className="text-lg font-medium">{user.name || user.email}</h3>
            <p className="text-sm text-gray-500">{user.email}</p>
          </div>
        </div>
        <div className="border-t pt-4 space-y-3">
          <div className="flex justify-between">
            <span className="text-gray-500">状态</span>
            <Badge variant={user.status === "active" ? "default" : "destructive"}>
              {user.status === "active" ? "正常" : "禁用"}
            </Badge>
          </div>
          <div className="flex justify-between">
            <span className="text-gray-500">角色</span>
            <div className="flex gap-1">
              {(user.roles || []).map((role: string) => (
                <Badge key={role} variant="outline">{role}</Badge>
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
