import { useState, useEffect } from 'react';
import ProductList from '../components/ProductList';

export default function Home() {
  const [products, setProducts] = useState([]);

  useEffect(() => {
    const fetchProducts = async () => {
      const res = await fetch('/api/products');
      const data = await res.json();
      setProducts(data);
    };
    fetchProducts();
    const interval = setInterval(fetchProducts, 60000);
    return () => clearInterval(interval);
  }, []);

  return (
    <div>
      <h1>Kith Restocks Monitor</h1>
      <ProductList products={products} />
    </div>
  );
}