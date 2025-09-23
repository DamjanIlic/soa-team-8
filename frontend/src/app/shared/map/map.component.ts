import { CommonModule } from '@angular/common';
import { AfterViewInit, Component, EventEmitter, Input, OnChanges, OnDestroy, Output, SimpleChanges } from '@angular/core';
import * as L from 'leaflet';
import 'leaflet-routing-machine';
import { Checkpoint } from './map.model';

const MAPBOX_TOKEN = 'pk.eyJ1IjoicHN3Z3J1cGEyIiwiYSI6ImNtMmc5OWlybTAwNHEya3F4emZrMDVoZGsifQ.aD0uouzJcAGE--8As0GFjg';

@Component({
  selector: 'xp-map',
  standalone: true,
  imports: [CommonModule],
  template: `<div id="map" style="height: 100%; width: 100%; border: 1px solid #ccc;"></div>`,
  styles: ['#map { width: 100%; height: 100%; }']
})
export class MapComponent implements AfterViewInit, OnChanges, OnDestroy {
  private map!: L.Map;
  private markers: L.Marker[] = [];
  private routingControl: any = null;

  @Input() addedCheckpointCollection: Checkpoint[] = [];
  @Output() locationSelected = new EventEmitter<{ lat: number, lng: number }>();
  @Output() checkpointRemoved = new EventEmitter<number>();
  @Output() routeDistanceKm = new EventEmitter<number>();

  ngAfterViewInit(): void {
    this.initMap();
  }

  ngOnChanges(changes: SimpleChanges): void {
    if (changes['addedCheckpointCollection'] && !changes['addedCheckpointCollection'].firstChange) {
      this.updateMarkersAndRoute();
    }
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
    });

    this.updateMarkersAndRoute();
  }

  private updateMarkersAndRoute(): void {
    // Ukloni stare markere i rutu
    this.markers.forEach(m => this.map.removeLayer(m));
    this.markers = [];
    if (this.routingControl) {
      this.map.removeControl(this.routingControl);
      this.routingControl = null;
    }

    // Dodaj markere
    this.addedCheckpointCollection.forEach((cp, i) => {
      const marker = L.marker([cp.latitude, cp.longitude], { title: cp.name })
        .addTo(this.map)
        .bindPopup(`<div style="width:200px"><strong>${cp.name}</strong><br>${cp.description || ''}<br>
        <button data-index="${i}">Remove</button></div>`);

      marker.on('popupopen', () => {
        const btn = document.querySelector(`button[data-index="${i}"]`);
        if (btn) btn.addEventListener('click', () => {
          this.checkpointRemoved.emit(i);
        });
      });

      this.markers.push(marker);
    });

    // Dodaj rutu ako ima >=2 checkpointa
    if (this.addedCheckpointCollection.length >= 2) {
      const waypoints = this.addedCheckpointCollection.map(cp => L.latLng(cp.latitude, cp.longitude));

      this.routingControl = (L as any).Routing.control({
        waypoints,
        router: (L as any).Routing.mapbox(MAPBOX_TOKEN, { profile: 'mapbox/driving' }),
        lineOptions: { styles: [{ color: 'blue', weight: 4 }] },
        addWaypoints: false,
        draggableWaypoints: false,
        fitSelectedRoutes: true,
        show: false
      }).addTo(this.map);

      this.routingControl.on('routesfound', (e: any) => {
        const distanceKm = e.routes[0].summary.totalDistance / 1000;
        this.routeDistanceKm.emit(distanceKm);
      });
    } else {
      this.routeDistanceKm.emit(0);
    }
  }

  ngOnDestroy(): void {
    if (this.map) this.map.remove();
  }
}
