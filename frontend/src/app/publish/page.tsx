"use client";

import { useCallback, useEffect, useState } from "react";
import { PublishTable } from "@/components/publish/publish-table";
import { api } from "@/lib/api";

export default function PublishPage() {
  const [items, setItems] = useState<any[]>([]);

  const fetchData = useCallback(async () => {
    const res = await api.publishTasks.list({ page: "1", pageSize: "50" });
    setItems(res.items || []);
  }, []);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  return (
    <div>
      <h2 className="text-2xl font-bold mb-6">发品任务</h2>
      <PublishTable items={items} />
    </div>
  );
}
