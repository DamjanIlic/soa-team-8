import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface FollowRequest {
  follower_id: string; 
  following_id: string;
}
//intercept:(

@Injectable({
  providedIn: 'root'
})
export class FollowService {

  private apiUrl = 'http://localhost:8000/api/follow';

  constructor(private http: HttpClient) {}

  // helper za ekstraktovanje userId iz JWT-a
  getCurrentUserId(): string | null {
    const token = localStorage.getItem('access_token'); // gde čuvaš JWT
    if (!token) return null;

    try {
      const payload = JSON.parse(atob(token.split('.')[1]));
      return payload.user_id; // ili kako se polje zove u JWT
    } catch (e) {
      console.error('Invalid token', e);
      return null;
    }
  }

  followUser(followingId: string): Observable<any> {
    const followerId = this.getCurrentUserId();
    if (!followerId) throw new Error("User not logged in");

    const request: FollowRequest = { follower_id: followerId, following_id: followingId };
    return this.http.post(`${this.apiUrl}/add`, request);
  }

  unfollowUser(followingId: string): Observable<any> {
    const followerId = this.getCurrentUserId();
    if (!followerId) throw new Error("User not logged in");

    const request: FollowRequest = { follower_id: followerId, following_id: followingId };
    return this.http.post(`${this.apiUrl}/remove`, request);
  }

  getFollowing(): Observable<any> {
    const userId = this.getCurrentUserId();
    if (!userId) throw new Error("User not logged in");

    return this.http.get(`${this.apiUrl}/following/${userId}`);
  }

  getFollowers(): Observable<any> {
    const userId = this.getCurrentUserId();
    if (!userId) throw new Error("User not logged in");

    return this.http.get(`${this.apiUrl}/followers/${userId}`);
  }

  getRecommendations(): Observable<any> {
    const userId = this.getCurrentUserId();
    if (!userId) throw new Error("User not logged in");

    return this.http.get(`${this.apiUrl}/recommendations/${userId}`);
  }
}
