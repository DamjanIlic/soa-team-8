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
}

