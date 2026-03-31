"use client";

import { useState } from "react";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { ProductEditDialog } from "./product-edit-dialog";
import { api } from "@/lib/api";

interface Product {
  ID: number;
  Name: string;
  Description: string;
  CostPrice: number;
  SellPrice: number;
  Status: string;
  CategoryID: string;
}

const statusMap: Record<string, { label: string; variant: "default" | "secondary" | "outline" }> = {
  draft: { label: "草稿", variant: "secondary" },
  ready: { label: "就绪", variant: "default" },
  published: { label: "已发布", variant: "outline" },
};

export function ProductTable({
  items,
  onRefresh,
}: {
  items: Product[];
  onRefresh: () => void;
}) {
  const [editProduct, setEditProduct] = useState<Product | null>(null);

  const handleMarkReady = async (id: number) => {
    try {
      await api.products.markReady(id);
      onRefresh();
    } catch (e: any) {
      alert(e.message);
    }
  };

  const handlePublish = async (id: number) => {
    try {
      await api.publishTasks.create({
        productId: id,
        targetPlatform: "pdd",
        categoryId: "cat-001",
        freightTemplate: "tpl-001",
      });
      alert("发品任务已创建，请到发品任务页查看");
      onRefresh();
    } catch (e: any) {
      alert(e.message);
    }
  };

  return (
    <>
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>商品名称</TableHead>
            <TableHead>成本价</TableHead>
            <TableHead>售价</TableHead>
            <TableHead>利润</TableHead>
            <TableHead>状态</TableHead>
            <TableHead>操作</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {items.map((item) => {
            const status = statusMap[item.Status] || statusMap.draft;
            const profit = item.SellPrice > 0 ? (item.SellPrice - item.CostPrice).toFixed(2) : "-";
            return (
              <TableRow key={item.ID}>
                <TableCell className="max-w-[250px] truncate">{item.Name}</TableCell>
                <TableCell>¥{item.CostPrice}</TableCell>
                <TableCell>{item.SellPrice > 0 ? `¥${item.SellPrice}` : "-"}</TableCell>
                <TableCell>{profit !== "-" ? `¥${profit}` : "-"}</TableCell>
                <TableCell>
                  <Badge variant={status.variant}>{status.label}</Badge>
                </TableCell>
                <TableCell className="space-x-2">
                  <Button size="sm" variant="outline" onClick={() => setEditProduct(item)}>
                    编辑
                  </Button>
                  {item.Status === "draft" && (
                    <Button size="sm" variant="outline" onClick={() => handleMarkReady(item.ID)}>
                      标记就绪
                    </Button>
                  )}
                  {item.Status === "ready" && (
                    <Button size="sm" onClick={() => handlePublish(item.ID)}>
                      发布到PDD
                    </Button>
                  )}
                </TableCell>
              </TableRow>
            );
          })}
          {items.length === 0 && (
            <TableRow>
              <TableCell colSpan={6} className="text-center text-gray-400 py-8">
                暂无商品，请先从货源管理页选品
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
      <ProductEditDialog
        product={editProduct}
        open={!!editProduct}
        onOpenChange={(open) => !open && setEditProduct(null)}
        onSuccess={onRefresh}
      />
    </>
  );
}
