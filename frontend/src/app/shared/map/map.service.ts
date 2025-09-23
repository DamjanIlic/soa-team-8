import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class MapService {
  constructor(private http: HttpClient) {}

  // opcionalno, ako želiš pretragu adresa
  search(street: string): Observable<any> {
    return this.http.get(`https://nominatim.openstreetmap.org/search?format=json&q=${street}`);
  }

  reverseSearch(lat: number, lon: number): Observable<any> {
    return this.http.get(`https://nominatim.openstreetmap.org/reverse?format=json&lat=${lat}&lon=${lon}`);
  }
}
