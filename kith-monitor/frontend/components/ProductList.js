import ProductCard from './ProductCard';

export default function ProductList({ products }) {
  return (
    <div className="product-list">
      {/* Message si aucun produit */}
      {products.length === 0 && (
        <div className="empty-message">
          Aucun produit restocké pour le moment. Vérifiez à nouveau plus tard !
        </div>
      )}

      {/* Liste des produits */}
      <div className="grid-container">
        {products.map((product) => (
          <ProductCard 
            key={product.id} 
            product={product} 
          />
        ))}
      </div>

      <style jsx>{`
        .product-list {
          max-width: 1200px;
          margin: 0 auto;
          padding: 20px;
        }
        .empty-message {
          text-align: center;
          font-size: 1.2rem;
          color: #666;
          margin: 40px 0;
        }
        .grid-container {
          display: grid;
          grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
          gap: 25px;
          padding: 20px 0;
        }
        @media (max-width: 768px) {
          .grid-container {
            grid-template-columns: 1fr;
          }
        }
      `}</style>
    </div>
  );
}