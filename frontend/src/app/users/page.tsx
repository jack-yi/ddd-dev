"use client";

import { useCallback, useEffect, useState } from "react";
import { UserTable } from "@/components/users/user-table";
import { api } from "@/lib/api";

export default function UsersPage() {
  const [items, setItems] = useState<any[]>([]);

  const fetchData = useCallback(async () => {
    const res = await api.users.list({ page: "1", pageSize: "50" });
    setItems(res.items || []);
  }, []);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  return (
    <div>
      <h2 className="text-2xl font-bold mb-6">用户管理</h2>
      <UserTable items={items} onRefresh={fetchData} />
    </div>
  );
}
