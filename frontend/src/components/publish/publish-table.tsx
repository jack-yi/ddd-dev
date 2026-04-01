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

interface PublishTask {
  ID: number;
  ProductID: number;
  TargetPlatform: string;
  PlatformProductID: string;
  Status: string;
  ErrorMessage: string;
  CreatedAt: string;
}

const statusMap: Record<
  string,
  { label: string; variant: "default" | "secondary" | "destructive" | "outline" }
> = {
  pending: { label: "待发布", variant: "secondary" },
  publishing: { label: "发布中", variant: "outline" },
  success: { label: "成功", variant: "default" },
  failed: { label: "失败", variant: "destructive" },
};

export function PublishTable({ items }: { items: PublishTask[] }) {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>任务ID</TableHead>
          <TableHead>商品ID</TableHead>
          <TableHead>目标平台</TableHead>
          <TableHead>平台商品ID</TableHead>
          <TableHead>状态</TableHead>
          <TableHead>错误信息</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {items.map((item) => {
          const status = statusMap[item.Status] || statusMap.pending;
          return (
            <TableRow key={item.ID}>
              <TableCell>{item.ID}</TableCell>
              <TableCell>{item.ProductID}</TableCell>
              <TableCell>{item.TargetPlatform}</TableCell>
              <TableCell>{item.PlatformProductID || "-"}</TableCell>
              <TableCell>
                <Badge variant={status.variant}>{status.label}</Badge>
              </TableCell>
              <TableCell className="max-w-[200px] truncate text-red-500">
                {item.ErrorMessage || "-"}
              </TableCell>
            </TableRow>
          );
        })}
        {items.length === 0 && (
          <TableRow>
            <TableCell colSpan={6} className="text-center text-gray-400 py-8">
              暂无发品任务
            </TableCell>
          </TableRow>
        )}
      </TableBody>
    </Table>
  );
}
