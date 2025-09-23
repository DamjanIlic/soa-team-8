export interface Checkpoint {
  id?: string;
  name: string;
  description: string;
  secret: string;
  latitude: number;
  longitude: number;
  image?: Image;
}

export interface Image {
  data: string;       // base64 string
  mimeType: string;   // npr. "image/png"
  uploadedAt: string; // ISO timestamp
}
