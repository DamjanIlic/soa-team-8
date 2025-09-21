import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface FollowRequest {
  follower_id: string; 
  following_id: string;
}
                              //add interceptor

@Injectable({
  providedIn: 'root'
})
export class FollowService {

  private apiUrl = 'http://localhost:8000/api/follow';

  constructor(private http: HttpClient) {}

  followUser(request: FollowRequest): Observable<any> {
    return this.http.post(`${this.apiUrl}/add`, request);
  }

  unfollowUser(request: FollowRequest): Observable<any> {
    return this.http.post(`${this.apiUrl}/remove`, request);
  }

  getFollowing(userId: string): Observable<any> {
    return this.http.get(`${this.apiUrl}/following/${userId}`);
  }

  getFollowers(userId: string): Observable<any> {
    return this.http.get(`${this.apiUrl}/followers/${userId}`);
  }

  getRecommendations(userId: string): Observable<any> {
    return this.http.get(`${this.apiUrl}/recommendations/${userId}`);
  }
}
