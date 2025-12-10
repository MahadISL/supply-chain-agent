"use client";
import { useEffect, useState } from 'react';
import { getProducts } from '../services/api'; // Path adjusts since services is in root

export default function Dashboard() {
  const [products, setProducts] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    getProducts()
      .then((data) => {
        setProducts(data);
        setLoading(false);
      })
      .catch((err) => {
        console.error(err);
        setLoading(false);
      });
  }, []);

  return (
    <div className="min-h-screen p-8">
      <header className="mb-8 flex justify-between items-center">
        <h1 className="text-3xl font-bold text-blue-800">Supply Chain Control Tower</h1>
        <a href="/agent" className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700">
          Open Agent Chat →
        </a>
      </header>

      <main>
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-xl font-semibold mb-4">Real-Time Inventory</h2>
          {loading ? (
            <p>Loading warehouse data...</p>
          ) : (
            <div className="overflow-x-auto">
              <table className="min-w-full text-left">
                <thead>
                  <tr className="border-b bg-gray-50">
                    <th className="p-3">ID</th>
                    <th className="p-3">Product Name</th>
                    <th className="p-3">Supplier</th>
                    <th className="p-3 text-right">Stock Level</th>
                    <th className="p-3">Status</th>
                  </tr>
                </thead>
                <tbody>
                  {products.map((p) => {
                    const isLow = p.stockQuantity < p.minStockLevel;
                    return (
                      <tr key={p.id} className="border-b hover:bg-gray-50">
                        <td className="p-3 text-gray-500">#{p.id}</td>
                        <td className="p-3 font-medium">{p.name}</td>
                        <td className="p-3 text-gray-600">{p.supplier.name}</td>
                        <td className="p-3 text-right font-mono">{p.stockQuantity}</td>
                        <td className="p-3">
                          {isLow ? (
                            <span className="bg-red-100 text-red-800 text-xs px-2 py-1 rounded-full font-bold">
                              LOW STOCK
                            </span>
                          ) : (
                            <span className="bg-green-100 text-green-800 text-xs px-2 py-1 rounded-full">
                              Healthy
                            </span>
                          )}
                        </td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </div>
          )}
        </div>
      </main>
    </div>
  );
}