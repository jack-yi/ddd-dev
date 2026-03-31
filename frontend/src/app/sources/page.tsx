"use client";

import { useCallback, useEffect, useState } from "react";
import { ImportDialog } from "@/components/sources/import-dialog";
import { SourceTable } from "@/components/sources/source-table";
import { Input } from "@/components/ui/input";
import { api } from "@/lib/api";

export default function SourcesPage() {
  const [items, setItems] = useState<any[]>([]);
  const [keyword, setKeyword] = useState("");

  const fetchData = useCallback(async () => {
    const params: Record<string, string> = { page: "1", pageSize: "50" };
    if (keyword) params.keyword = keyword;
    const res = await api.sourceItems.list(params);
    setItems(res.items || []);
  }, [keyword]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-2xl font-bold">货源管理</h2>
        <ImportDialog onSuccess={fetchData} />
      </div>
      <div className="mb-4">
        <Input
          placeholder="搜索货源标题..."
          value={keyword}
          onChange={(e) => setKeyword(e.target.value)}
          className="max-w-sm"
        />
      </div>
      <SourceTable items={items} onRefresh={fetchData} />
    </div>
  );
}
