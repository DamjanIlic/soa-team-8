import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, map } from 'rxjs';
import { Blog } from '../models/blog.model';
import { Comment } from '../models/comment.model';

@Injectable({
  providedIn: 'root'
})
export class BlogService {
  private apiUrl = 'http://localhost:8000/api/blogs';

  constructor(private http: HttpClient) {}

  // ================== BLOG METHODS ==================

  getAll(): Observable<Blog[]> {
    return this.http.get<Blog[]>(this.apiUrl);
  }

  getById(id: string): Observable<Blog> {
    return this.http.get<Blog>(`${this.apiUrl}/${id}`);
  }

  create(blog: Omit<Blog, 'id' | 'created_at' | 'updated_at' | 'likes'>): Observable<Blog> {
    return this.http.post<Blog>(this.apiUrl, blog);
  }

  like(blogId: string): Observable<number> {
    return this.http.post<{ likes: number }>(
      `${this.apiUrl}/${blogId}/like`,
      {}
    ).pipe(map(res => res.likes));
  }

  unlike(blogId: string): Observable<number> {
    return this.http.post<{ likes: number }>(
      `${this.apiUrl}/${blogId}/unlike`,
      {}
    ).pipe(map(res => res.likes));
  }

    // helper za ekstraktovanje userId iz JWT
  private getCurrentUserId(): string | null {
    const token = localStorage.getItem('access_token');
    if (!token) return null;

    try {
      const payload = JSON.parse(atob(token.split('.')[1]));
      return payload.user_id;
    } catch (e) {
      console.error('Invalid token', e);
      return null;
    }
  }

  // Funkcija koja vraća blogove za current usera (onog ko je ulogovan)
  getForCurrentUser(): Observable<Blog[]> {
    const userId = this.getCurrentUserId();
    if (!userId) {
      throw new Error('No current user ID found in token');
    }

    return this.http.get<Blog[]>(`http://localhost:8000/api/blogs/user/${userId}`);
  }

  // ================== COMMENT METHODS ==================

  getComments(blogId: string): Observable<Comment[]> {
    return this.http.get<Comment[]>(`${this.apiUrl}/${blogId}/comments`);
  }

  createComment(blogId: string, comment: Comment): Observable<Comment> {
    return this.http.post<Comment>(
      `${this.apiUrl}/${blogId}/comments`,
      comment
    );
  }
}
