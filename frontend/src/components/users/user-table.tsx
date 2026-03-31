"use client";

import {
  Table, TableBody, TableCell, TableHead, TableHeader, TableRow,
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { api } from "@/lib/api";
import { getUser } from "@/lib/auth";

interface User {
  ID: number;
  Email: string;
  Name: string;
  Status: string;
}

const roleOptions = ["super_admin", "admin", "operator", "viewer"];

export function UserTable({ items, onRefresh }: { items: User[]; onRefresh: () => void }) {
  const currentUser = getUser();
  const isSuperAdmin = currentUser?.roles?.includes("super_admin");

  const handleToggleStatus = async (id: number, currentStatus: string) => {
    const newStatus = currentStatus === "active" ? "disabled" : "active";
    await api.users.updateStatus(id, newStatus);
    onRefresh();
  };

  const handleAssignRole = async (id: number, roleName: string) => {
    await api.users.assignRole(id, roleName);
    onRefresh();
  };

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>邮箱</TableHead>
          <TableHead>名称</TableHead>
          <TableHead>状态</TableHead>
          <TableHead>操作</TableHead>
          {isSuperAdmin && <TableHead>角色分配</TableHead>}
        </TableRow>
      </TableHeader>
      <TableBody>
        {items.map((item) => (
          <TableRow key={item.ID}>
            <TableCell>{item.Email}</TableCell>
            <TableCell>{item.Name}</TableCell>
            <TableCell>
              <Badge variant={item.Status === "active" ? "default" : "destructive"}>
                {item.Status === "active" ? "正常" : "禁用"}
              </Badge>
            </TableCell>
            <TableCell>
              <Button size="sm" variant="outline" onClick={() => handleToggleStatus(item.ID, item.Status)}>
                {item.Status === "active" ? "禁用" : "启用"}
              </Button>
            </TableCell>
            {isSuperAdmin && (
              <TableCell>
                <select
                  className="border rounded px-2 py-1 text-sm"
                  onChange={(e) => handleAssignRole(item.ID, e.target.value)}
                  defaultValue=""
                >
                  <option value="" disabled>选择角色</option>
                  {roleOptions.map((r) => (
                    <option key={r} value={r}>{r}</option>
                  ))}
                </select>
              </TableCell>
            )}
          </TableRow>
        ))}
        {items.length === 0 && (
          <TableRow>
            <TableCell colSpan={5} className="text-center text-gray-400 py-8">暂无用户</TableCell>
          </TableRow>
        )}
      </TableBody>
    </Table>
  );
}
