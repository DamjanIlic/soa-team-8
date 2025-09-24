import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable, map } from 'rxjs';
import { Blog, BlogRequest, BlogResponse } from '../models/blog.model';
import { Comment } from '../models/comment.model';

@Injectable({
  providedIn: 'root'
})
export class BlogService {
  private apiUrl = 'http://localhost:8000/api/blogs';

  constructor(private http: HttpClient) {}

  private getAuthHeaders(): HttpHeaders {
    const token = localStorage.getItem('access_token');
    return new HttpHeaders({
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    });
  }

  // ================== BLOG METHODS ==================

  getAll(): Observable<Blog[]> {
    return this.http.get<Blog[]>(this.apiUrl, { headers: this.getAuthHeaders() });
  }

  getMyBlogs(): Observable<Blog[]> {
  return this.http.get<Blog[]>(`${this.apiUrl}/my`, { headers: this.getAuthHeaders() });
  }

  getById(id: string): Observable<Blog> {
    return this.http.get<Blog>(`${this.apiUrl}/${id}`, { headers: this.getAuthHeaders() });
  }

  create(blog: BlogRequest): Observable<BlogResponse> {
    return this.http.post<BlogResponse>(this.apiUrl, blog, { headers: this.getAuthHeaders() });
  }

  like(blogId: string): Observable<number> {
    return this.http.post<{ likes: number }>(
      `${this.apiUrl}/${blogId}/like`,
      {}, 
      { headers: this.getAuthHeaders() }
    ).pipe(map(res => res.likes));
  }

  unlike(blogId: string): Observable<number> {
    return this.http.post<{ likes: number }>(
      `${this.apiUrl}/${blogId}/unlike`,
      {}, 
      { headers: this.getAuthHeaders() }
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

    return this.http.get<Blog[]>(`${this.apiUrl}/user/${userId}`, { headers: this.getAuthHeaders() });
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
