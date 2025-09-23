export interface Blog {
  id: string;
  UserID?: string;
  title: string;
  content: string;
  image_url: string;
  created_at: string;
  updated_at: string;
  likes: number;
}