import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable, map } from 'rxjs';
import { Blog } from '../models/blog.model';
import { Comment } from '../models/comment.model';
@Injectable({
  providedIn: 'root'
})
export class BlogService {
  private apiUrl = 'http://localhost:8000/api/blogs';

  constructor(private http: HttpClient) {}

  getAll(): Observable<Blog[]> {
    return this.http.get<Blog[]>(this.apiUrl);
  }

  getById(id: string): Observable<Blog> {
    return this.http.get<Blog>(`${this.apiUrl}/${id}`);
  }

  create(blog: Omit<Blog, 'id' | 'created_at' | 'updated_at' | 'likes'>, token: string): Observable<Blog> {
    return this.http.post<Blog>(this.apiUrl, blog, {
      headers: { Authorization: `Bearer ${token}` }
    });
  }

  like(blogId: string, token: string): Observable<number> {
    return this.http.post<{ likes: number }>(
      `${this.apiUrl}/${blogId}/like`,
      {},
      { headers: { Authorization: `Bearer ${token}` } }
    ).pipe(map(res => res.likes));
  }

  unlike(blogId: string, token: string): Observable<number> {
    return this.http.post<{ likes: number }>(
      `${this.apiUrl}/${blogId}/unlike`,
      {},
      { headers: { Authorization: `Bearer ${token}` } }
    ).pipe(map(res => res.likes));
  }

  // ================== COMMENT METHODS ==================

  getComments(blogId: string): Observable<Comment[]> {
    return this.http.get<Comment[]>(`${this.apiUrl}/${blogId}/comments`);
  }

  createComment(blogId: string, comment: Comment, token: string): Observable<Comment> {
    const headers = new HttpHeaders({
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    });

    return this.http.post<Comment>(
      `${this.apiUrl}/${blogId}/comments`,
      comment,
      { headers }
    );
  }
}