import { District } from "./constants.ts";

type DistrictInfo = {
  hospital_obs: number;
  home_obs: number;
  total_obs: number;
  hospital_today: number;
  confirmed: number;
  recovered: number;
  deceased: number;
  active: number;
};

type HistoryInfo = Record<District, DistrictInfo>;

type History = {
  summary: HistoryInfo;
  delta: HistoryInfo;
  date: string;
};

type Histories = {
  histories: History[];
  last_updated: string;
};

type LatestHistory = {
  summary: HistoryInfo;
  delta: HistoryInfo;
  last_updated: string;
};

type Summary = {
  summary: DistrictInfo;
  delta: DistrictInfo;
  last_updated: string;
};

type TestReport = {
  date: string;
  total: number;
  today: number;
  positive: number;
  today_positive: number;
};

type TestReports = {
  reports: TestReport[];
  last_updated: string;
};

type Hotspots = {
  district: string;
  lsgd: string;
  wards: string;
};

type HotspotsHistory = {
  hotspots: Hotspots[];
  date: string;
};

type HotspotsHistories = {
  histories: HotspotsHistory[];
  last_updated: string;
};

type LatestHotspotsHistory = {
  hotspots: Hotspots[];
  last_updated: string;
};

type VacDistrictInfo = {
  cvc_public: number;
  cvc_private: number;
  cvc_total: number;
  hcw_dose1: number;
  hcw_dose2: number;
  flw_other_dose1: number;
  flw_other_dose2: number;
  flw_polling_dose1: number;
  flw_polling_dose2: number;
  age_appropriate_dose1: number;
  age_appropriate_dose2: number;
  tot_vaccinations: number;
  tot_person_vaccinations: number;
  second_dose: number;
};

type VacHistoryInfo = Record<District, VacDistrictInfo>;

type VacHistory = {
  summary: VacHistoryInfo;
  delta: VacHistoryInfo;
  date: string;
};

type VacHistories = {
  histories: VacHistory[];
  last_updated: string;
};

type VacLatestHistory = {
  summary: VacHistoryInfo;
  delta: VacHistoryInfo;
  last_updated: string;
};

type VacSummary = {
  summary: VacDistrictInfo;
  delta: VacDistrictInfo;
  last_updated: string;
};
