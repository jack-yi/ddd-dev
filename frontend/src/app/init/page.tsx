"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { api } from "@/lib/api";
import { getUser, setUser } from "@/lib/auth";

export default function InitPage() {
  const router = useRouter();
  const [needInit, setNeedInit] = useState<boolean | null>(null);
  const [loading, setLoading] = useState(false);
  const user = getUser();

  useEffect(() => {
    api.auth.checkInit().then((res) => setNeedInit(res.needInit));
  }, []);

  const handleInit = async () => {
    setLoading(true);
    try {
      await api.auth.initSuperAdmin();
      const updatedUser = await api.auth.me();
      setUser(updatedUser);
      alert("超级管理员初始化成功！");
      router.replace("/sources");
    } catch (e: any) {
      alert(e.message);
    } finally {
      setLoading(false);
    }
  };

  if (needInit === null) return <div className="p-8">检查中...</div>;
  if (!needInit) {
    router.replace("/sources");
    return null;
  }

  return (
    <div className="max-w-md mx-auto mt-20 p-8 bg-white rounded-xl shadow-md space-y-6">
      <h2 className="text-xl font-bold">初始化超级管理员</h2>
      <p className="text-gray-500">
        系统首次使用，需要初始化超级管理员。当前登录用户
        <strong> {user?.email} </strong>
        将被设为超级管理员。
      </p>
      <Button onClick={handleInit} disabled={loading} className="w-full">
        {loading ? "初始化中..." : "确认初始化"}
      </Button>
    </div>
  );
}
