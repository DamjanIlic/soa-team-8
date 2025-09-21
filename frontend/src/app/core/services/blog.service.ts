import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, map } from 'rxjs';
import { Blog } from '../models/blog.model';

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
}
