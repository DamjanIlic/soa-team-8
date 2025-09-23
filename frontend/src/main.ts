import { bootstrapApplication } from '@angular/platform-browser';
import { AppComponent } from './app/app.component';
import { appConfig } from './app/app.config';

// main.ts ili styles.ts
import 'leaflet/dist/leaflet.css';

bootstrapApplication(AppComponent, appConfig)
  .catch((err) => console.error(err));
