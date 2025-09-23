import { CommonModule } from '@angular/common';
import { AfterViewInit, Component, EventEmitter, Input, OnChanges, OnDestroy, Output, SimpleChanges } from '@angular/core';
import * as L from 'leaflet';

import { Checkpoint } from './map.model';

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
  private polyline: L.Polyline | null = null;

  @Input() addedCheckpointCollection: Checkpoint[] = [];
  @Input() clearMarkersTrigger: boolean = false;
  @Output() locationSelected = new EventEmitter<{ lat: number; lng: number }>();
  @Output() checkpointRemoved = new EventEmitter<number>();

  ngAfterViewInit(): void {
    this.initMap();
  }

  ngOnChanges(changes: SimpleChanges): void {
    if (changes['addedCheckpointCollection'] && !changes['addedCheckpointCollection'].firstChange) {
      this.updateMarkersAndLine();
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

    this.updateMarkersAndLine();
  }

  private updateMarkersAndLine(): void {
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

    // Crtaj poliliniju kroz sve checkpoint-e
    if (this.addedCheckpointCollection.length >= 2) {
      const latlngs = this.addedCheckpointCollection.map(cp => [cp.latitude, cp.longitude] as [number, number]);
      this.polyline = L.polyline(latlngs, { color: 'blue' }).addTo(this.map);
    }
  }

  clearMarkers(): void {
    this.markers.forEach(m => this.map.removeLayer(m));
    this.markers = [];
    if (this.polyline) {
      this.map.removeLayer(this.polyline);
      this.polyline = null;
    }
  }

  ngOnDestroy(): void {
    if (this.map) this.map.remove();
  }
}
