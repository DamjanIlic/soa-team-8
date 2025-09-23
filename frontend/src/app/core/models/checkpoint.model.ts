export interface Checkpoint {
  id?: string;          // opciono, generiše backend
  name: string;         // naziv checkpoint-a
  description: string;  // opis checkpoint-a
  latitude: number;     // geografska širina
  longitude: number;    // geografska dužina
  image_url?: string;   // URL slike, opciono
}
