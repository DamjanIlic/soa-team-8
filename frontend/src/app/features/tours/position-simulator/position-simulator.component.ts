// position-simulator.component.ts
import { Component, OnInit } from '@angular/core';
import * as L from 'leaflet';

@Component({
  selector: 'app-position-simulator',
  templateUrl: './position-simulator.component.html'
})
export class PositionSimulatorComponent implements OnInit {
  private map!: L.Map;
  private userMarker!: L.Marker;

  ngOnInit(): void {
    this.initMap();
  }

  private initMap(): void {
    // Početna pozicija korisnika
    const initialPosition: L.LatLngExpression = [44.8176, 20.4569];

    // Inicijalizacija mape
    this.map = L.map('map').setView(initialPosition, 15);

    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '&copy; OpenStreetMap contributors'
    }).addTo(this.map);

    // Marker korisnika
    this.userMarker = L.marker(initialPosition)
      .addTo(this.map)
      .bindPopup("User is here")
      .openPopup();

    // Klik na mapu pomera korisnika
    this.map.on('click', (e: L.LeafletMouseEvent) => {
      this.userMarker.setLatLng(e.latlng);
      this.userMarker.getPopup()?.setContent(`User is here: ${e.latlng.lat.toFixed(5)}, ${e.latlng.lng.toFixed(5)}`);
      this.userMarker.openPopup();
    });
  }
}
