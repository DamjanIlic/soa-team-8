export interface Tour {
  id: string;
  author_id: string;
  name: string;
  description: string;
  difficulty: string;
  tags: string;
  status: 'draft' | 'published' | 'archived';
  price: number;
  distance_km: number;
  durations: Duration[];
  checkpoints: Checkpoint[];   // <-- dodato
  published_at?: string;
  archived_at?: string;
  created_at: string;
  updated_at: string;
}


export interface Duration {
  id?: string;
  tour_id?: string;
  transport: string;
  duration: number;
}

export interface TourRequest {
  name: string;
  description: string;
  difficulty: string;
  tags: string;
}

export interface TourResponse {
  id: string;
  author_id: string;
  name: string;
  description: string;
  difficulty: string;
  tags: string;
  status: 'draft' | 'published' | 'archived';
  price: number;  
  distance_km: number;
  durations: Duration[];
  created_at: string;
  updated_at: string;
  published_at?: string;
  archived_at?: string;
  checkpoints: Checkpoint[];
}

export interface Checkpoint {
  id?: string;          // opciono, generiše backend
  name: string;         // naziv checkpoint-a
  description: string;  // opis checkpoint-a
  latitude: number;     // geografska širina
  longitude: number;    // geografska dužina
  image_url?: string;   // URL slike, opciono
}
