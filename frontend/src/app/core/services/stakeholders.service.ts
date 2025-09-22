import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class StakeholdersService {
  private apiUrl = 'http://localhost:8000/api/stakeholders';

  constructor(private http: HttpClient) {}

  getProfile(): Observable<any> {
    const token = localStorage.getItem('access_token'); // JWT token
    if (!token) {
        throw new Error('No token found');
    }
    const headers = new HttpHeaders({
      'Authorization': `Bearer ${token}`
    });
    return this.http.get(`${this.apiUrl}/profile`, { headers });
  }


    // ================== ADMIN METHODS ==================

  // GET /stakeholders/admin/all
  getAllUsers(): Observable<any> {
    return this.http.get(`${this.apiUrl}/admin/all`);
  }

  // PUT /stakeholders/admin/users/{id}/block
  blockUser(userId: string): Observable<any> {
    return this.http.put(`${this.apiUrl}/admin/users/${userId}/block`, {});
  }

  unblockUser(userId: string): Observable<any> {
    return this.http.put(`${this.apiUrl}/admin/users/${userId}/unblock`, {});
  }

  updateProfile(profileData: any): Observable<any> {
    const token = localStorage.getItem('access_token');
    if (!token) {
      throw new Error('No token found');
    }
    const headers = new HttpHeaders({
      'Authorization': `Bearer ${token}`
    });
    return this.http.put(`${this.apiUrl}/profile`, profileData, { headers });
  }
}