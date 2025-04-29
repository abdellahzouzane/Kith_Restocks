export default function ProductCard({ product }) {
    return (
      <div className="product-card">
        <img src={product.image} alt={product.name} />
        <h3>{product.name}</h3>
        <p>Prix: {product.price} €</p>
        <a href={product.url} target="_blank">Voir sur Kith</a>
        <p>Détecté le: {new Date(product.detected_at).toLocaleString()}</p>
      </div>
    );
  }