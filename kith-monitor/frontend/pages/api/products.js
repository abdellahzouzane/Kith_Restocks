import { Pool } from 'pg';

const pool = new Pool({
  user: 'abdellah_user',
  host: '127.0.0.1',
  database: 'db_user',
  password: 'password',
  port: 5432,
});

export default async function handler(req, res) {
  try {
    const { rows } = await pool.query(`
      SELECT id, name, url, price, image, in_stock, detected_at
      FROM products
      WHERE in_stock = true
      ORDER BY detected_at DESC
    `);
    res.status(200).json(rows);
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
}