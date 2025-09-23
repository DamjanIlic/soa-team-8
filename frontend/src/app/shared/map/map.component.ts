import { CommonModule } from '@angular/common';
import { AfterViewInit, Component, EventEmitter, Input, OnChanges, OnDestroy, Output, SimpleChanges } from '@angular/core';
import * as L from 'leaflet';
import 'leaflet-routing-machine';
import { Checkpoint } from './map.model';

// TypeScript deklaracije za L.Routing
declare module 'leaflet' {
  namespace Routing {
    class Control extends L.Control {
      constructor(options?: any);
      on(type: string, fn: (e: any) => void): this;
      getPlan(): any;
      setWaypoints(waypoints: L.LatLng[]): void;
    }
    function control(options?: any): Control;
    function osrmv1(options?: any): any;
  }

  interface Map {
    routingControl?: Routing.Control;
  }
}

@Component({
  selector: 'xp-map',
  standalone: true,
  imports: [CommonModule],
  template: `<div id="map" style="height: 100%; width: 100%; border: 1px solid #ccc;"></div>`,
  styles: ['#map { width: 100%; height: 100%; }']
})
export class MapComponent implements AfterViewInit, OnChanges, OnDestroy {
  private map!: L.Map;
  private markers: L.Layer[] = [];
  private routingControl: L.Routing.Control | null = null;

  @Input() addedCheckpointCollection: Checkpoint[] = [];
  @Input() clearMarkersTrigger: boolean = false;
  @Output() locationSelected = new EventEmitter<{ lat: number; lng: number }>();
  @Output() checkpointRemoved = new EventEmitter<number>();
  @Output() routeDistanceKm = new EventEmitter<number>();

  ngAfterViewInit(): void {
    this.initMap();
  }

  ngOnChanges(changes: SimpleChanges): void {
    if (changes['addedCheckpointCollection'] && !changes['addedCheckpointCollection'].firstChange) {
      this.updateMarkersAndRoute();
    }
    if (changes['clearMarkersTrigger'] && changes['clearMarkersTrigger'].currentValue) {
      this.clearMarkers();
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
    this.clearMarkers();

    // Dodaj markere
    this.addedCheckpointCollection.forEach((cp, i) => {
      const marker = L.circleMarker([cp.latitude, cp.longitude], {
        radius: 8,
        color: 'red',
        fillColor: 'red',
        fillOpacity: 1
      }).addTo(this.map)
        .bindPopup(`${cp.name || 'Checkpoint'} <button type="button" data-index="${i}">Remove</button>`);

      marker.on('popupopen', () => {
        const btn = document.querySelector(`button[data-index="${i}"]`);
        if (btn) btn.addEventListener('click', () => {
          this.checkpointRemoved.emit(i);
        });
      });

      this.markers.push(marker);
    });

    // Crtaj rutu ako ima >= 2 checkpoint-a
    if (this.addedCheckpointCollection.length >= 2) {
      const waypoints = this.addedCheckpointCollection.map(cp => L.latLng(cp.latitude, cp.longitude));

      if (this.routingControl) this.map.removeControl(this.routingControl);

      this.routingControl = L.Routing.control({
        waypoints,
        router: (L.Routing as any).osrmv1({ serviceUrl: 'https://router.project-osrm.org/route/v1' }),
        lineOptions: { styles: [{ color: 'blue', weight: 4 }] },
        addWaypoints: false,
        draggableWaypoints: false,
        fitSelectedRoutes: true,
        show: false
      }).addTo(this.map);

      this.routingControl.on('routesfound', (e: any) => {
        const route = e.routes[0];
        const distanceKm = route.summary.totalDistance / 1000;
        this.routeDistanceKm.emit(distanceKm);
      });
    } else {
      this.routeDistanceKm.emit(0);
    }
  }

  clearMarkers(): void {
    this.markers.forEach(m => this.map.removeLayer(m));
    this.markers = [];
    if (this.routingControl) {
      this.map.removeControl(this.routingControl);
      this.routingControl = null;
    }
  }

  ngOnDestroy(): void {
    if (this.map) this.map.remove();
  }
}
