"use client";

import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { api } from "@/lib/api";

interface Product {
  ID: number;
  Name: string;
  Description: string;
  CostPrice: number;
  SellPrice: number;
}

export function ProductEditDialog({
  product,
  open,
  onOpenChange,
  onSuccess,
}: {
  product: Product | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess: () => void;
}) {
  const [name, setName] = useState("");
  const [sellPrice, setSellPrice] = useState("");
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (product) {
      setName(product.Name || "");
      setSellPrice(String(product.SellPrice || ""));
    }
  }, [product]);

  const handleSave = async () => {
    if (!product) return;
    setLoading(true);
    try {
      await api.products.update(product.ID, {
        name,
        sellPrice: parseFloat(sellPrice),
      });
      onOpenChange(false);
      onSuccess();
    } catch (e: any) {
      alert(e.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>编辑商品</DialogTitle>
        </DialogHeader>
        <div className="space-y-4 pt-4">
          <div>
            <Label>商品名称</Label>
            <Input value={name} onChange={(e) => setName(e.target.value)} />
          </div>
          <div>
            <Label>成本价</Label>
            <Input value={product?.CostPrice || 0} disabled />
          </div>
          <div>
            <Label>售价</Label>
            <Input
              type="number"
              value={sellPrice}
              onChange={(e) => setSellPrice(e.target.value)}
            />
          </div>
          <Button onClick={handleSave} disabled={loading} className="w-full">
            {loading ? "保存中..." : "保存"}
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
