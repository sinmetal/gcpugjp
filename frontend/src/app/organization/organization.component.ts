import { Component, OnInit } from '@angular/core';
import {Observable} from "rxjs/Observable";
import {OrganizationListAPIResponse} from "../shared/models/organization";
import {ActivatedRoute} from "@angular/router";
import {OrganizationService} from "../services/organization.service";

@Component({
  selector: 'app-organization',
  templateUrl: './organization.component.html',
  styleUrls: ['./organization.component.scss']
})
export class OrganizationComponent implements OnInit {
  public organizationAPIResObservable: Observable<OrganizationListAPIResponse>;

  constructor(
      private route: ActivatedRoute,
      private organizationService: OrganizationService
  ) { }

  ngOnInit() {
    this.organizationAPIResObservable = this.route.params.switchMap(() => {
      return this.organizationService.list();
    });
  }

}
