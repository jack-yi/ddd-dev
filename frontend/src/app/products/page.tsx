"use client";

import { useCallback, useEffect, useState } from "react";
import { ProductTable } from "@/components/products/product-table";
import { api } from "@/lib/api";

export default function ProductsPage() {
  const [items, setItems] = useState<any[]>([]);

  const fetchData = useCallback(async () => {
    const res = await api.products.list({ page: "1", pageSize: "50" });
    setItems(res.items || []);
  }, []);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  return (
    <div>
      <h2 className="text-2xl font-bold mb-6">商品管理</h2>
      <ProductTable items={items} onRefresh={fetchData} />
    </div>
  );
}
