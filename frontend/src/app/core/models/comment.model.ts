export interface Comment {
  id?: string;
  blogId: string;
  user_id?: string;
  username?: string;
  name?: string;
  surname?: string;
  avatarUrl?: string;
  motto?: string;
  text: string;
  created_at?: string;
}