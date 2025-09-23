import { AfterViewInit, Component, EventEmitter, Input, OnDestroy, Output, SimpleChanges } from '@angular/core';
import * as L from 'leaflet';
import 'leaflet-routing-machine';

import { Checkpoint } from './map.model';

@Component({
  selector: 'xp-map',
  template: `<div id="map" style="height: 500px;"></div>`,
  styles: ['#map { width: 100%; height: 100%; }']
})
export class MapComponent implements AfterViewInit, OnDestroy {
  map!: L.Map;
  private markers: L.Marker[] = [];
  private routingControl: L.Routing.Control | null = null;

  @Input() addedCheckpointCollection: Checkpoint[] = [];
  @Input() clearMarkersTrigger: boolean = false;
  @Output() locationSelected = new EventEmitter<{ lat: number; lng: number }>();
  @Output() checkpointRemoved = new EventEmitter<number>();

  ngAfterViewInit(): void {
    this.initMap();
  }

  private initMap(): void {
    this.map = L.map('map', { center: [45.2396, 19.8227], zoom: 13 });

    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      maxZoom: 18,
      attribution: '&copy; OpenStreetMap contributors'
    }).addTo(this.map);

    this.map.on('click', (e: any) => {
      const { lat, lng } = e.latlng;
      this.locationSelected.emit({ lat, lng });
      const marker = L.marker([lat, lng]).addTo(this.map);
      this.markers.push(marker);
    });

    this.addRouteToMap();
  }

  ngOnChanges(changes: SimpleChanges): void {
    if (changes['addedCheckpointCollection'] && !changes['addedCheckpointCollection'].firstChange) {
      this.clearMarkers();
      this.addRouteToMap();
    }
    if (changes['clearMarkersTrigger'] && changes['clearMarkersTrigger'].currentValue) {
      this.clearMarkers();
    }
  }

  private addRouteToMap(): void {
    if (this.addedCheckpointCollection.length < 1) return;

    const waypoints = this.addedCheckpointCollection.map(cp => L.latLng(cp.latitude, cp.longitude));

    waypoints.forEach((latlng, i) => {
      const marker = L.marker(latlng).addTo(this.map)
        .bindPopup(`${this.addedCheckpointCollection[i].name} <button type="button" data-index="${i}">Remove</button>`);
      marker.on('popupopen', () => {
        const btn = document.querySelector(`button[data-index="${i}"]`);
        if (btn) btn.addEventListener('click', () => {
          this.checkpointRemoved.emit(i);
          this.map.removeLayer(marker);
        });
      });
      this.markers.push(marker);
    });

    if (waypoints.length >= 2) {
      const plan = new L.Routing.Plan(waypoints, { draggableWaypoints: false, createMarker: () => false });
      this.routingControl = L.Routing.control({ plan, routeWhileDragging: false, addWaypoints: false }).addTo(this.map);
    }
  }

  clearMarkers(): void {
    this.markers.forEach(m => this.map.removeLayer(m));
    this.markers = [];
    if (this.routingControl) { this.map.removeControl(this.routingControl); this.routingControl = null; }
  }

  ngOnDestroy(): void {
    if (this.map) this.map.remove();
  }
}
