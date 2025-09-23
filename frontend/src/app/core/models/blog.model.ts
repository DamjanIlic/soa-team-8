export interface Blog {
  id: string;
  user_id?: string;
  username?: string;
  title: string;
  content: string;
  image_url: string;
  created_at: string;
  updated_at: string;
  likes: number;
}