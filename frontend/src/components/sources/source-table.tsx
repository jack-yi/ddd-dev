"use client";

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
import { api } from "@/lib/api";

interface SourceItem {
  ID: number;
  Platform: string;
  Title: string;
  PriceMin: number;
  PriceMax: number;
  Supplier: { Name: string; Rating: number; Region: string };
  Category: string;
  SalesVolume: number;
  Status: string;
}

const statusMap: Record<string, { label: string; variant: "default" | "secondary" | "destructive" }> = {
  new: { label: "新导入", variant: "secondary" },
  selected: { label: "已选品", variant: "default" },
  ignored: { label: "已忽略", variant: "destructive" },
};

export function SourceTable({
  items,
  onRefresh,
}: {
  items: SourceItem[];
  onRefresh: () => void;
}) {
  const handleSelect = async (id: number) => {
    await api.sourceItems.updateStatus(id, "selected");
    onRefresh();
  };

  const handleCreateProduct = async (id: number) => {
    await api.products.createFromSource(id);
    alert("商品创建成功，请到商品管理页编辑");
    onRefresh();
  };

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>标题</TableHead>
          <TableHead>平台</TableHead>
          <TableHead>价格区间</TableHead>
          <TableHead>供应商</TableHead>
          <TableHead>销量</TableHead>
          <TableHead>状态</TableHead>
          <TableHead>操作</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {items.map((item) => {
          const status = statusMap[item.Status] || statusMap.new;
          return (
            <TableRow key={item.ID}>
              <TableCell className="max-w-[200px] truncate">{item.Title}</TableCell>
              <TableCell>{item.Platform}</TableCell>
              <TableCell>
                ¥{item.PriceMin} - ¥{item.PriceMax}
              </TableCell>
              <TableCell>
                {item.Supplier.Name} ({item.Supplier.Rating})
              </TableCell>
              <TableCell>{item.SalesVolume}</TableCell>
              <TableCell>
                <Badge variant={status.variant}>{status.label}</Badge>
              </TableCell>
              <TableCell className="space-x-2">
                {item.Status === "new" && (
                  <Button size="sm" variant="outline" onClick={() => handleSelect(item.ID)}>
                    选品
                  </Button>
                )}
                {item.Status === "selected" && (
                  <Button size="sm" onClick={() => handleCreateProduct(item.ID)}>
                    创建商品
                  </Button>
                )}
              </TableCell>
            </TableRow>
          );
        })}
        {items.length === 0 && (
          <TableRow>
            <TableCell colSpan={7} className="text-center text-gray-400 py-8">
              暂无货源，点击「导入货源」开始
            </TableCell>
          </TableRow>
        )}
      </TableBody>
    </Table>
  );
}
