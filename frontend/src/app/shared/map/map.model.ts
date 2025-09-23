export interface Checkpoint {
  id?: string;
  name: string;
  description?: string;
  latitude: number;
  longitude: number;
  image?: { mimeType: string; data: string };
}
