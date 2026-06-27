import {
  AfterViewInit,
  Component,
  DestroyRef,
  ElementRef,
  OnInit,
  ViewChild,
  inject,
} from '@angular/core';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { MatButtonModule } from '@angular/material/button';
import { MatMenuModule } from '@angular/material/menu';
import { max, sum } from 'd3-array';
import { axisBottom, axisLeft } from 'd3-axis';
import { format } from 'd3-format';
import { scaleBand, scaleLinear } from 'd3-scale';
import { select, type Selection } from 'd3-selection';
import {
  arc,
  curveMonotoneX,
  line,
  pie,
  stack,
  type PieArcDatum,
} from 'd3-shape';
import { forkJoin } from 'rxjs';
import { ReportService } from './services';
import { ReportChartResponse } from './types';

type ChartDatum = { label: string; value: number };
type NetChangeDatum = { month: string; value: number };
type TierKey = 'Classic' | 'Standard' | 'Premium' | 'Credited';
type TierHistoryDatum = Record<TierKey, number> & { month: string };

const TIER_KEYS: TierKey[] = ['Classic', 'Standard', 'Premium', 'Credited'];
const TIER_COLORS = new Map<TierKey, string>([
  ['Classic', '#6f6cab'],
  ['Standard', '#e9437a'],
  ['Premium', '#36a269'],
  ['Credited', '#7c8796'],
]);

@Component({
  selector: 'md-report',
  imports: [MatButtonModule, MatMenuModule],
  providers: [ReportService],
  templateUrl: './report.component.html',
  styleUrl: './report.component.scss',
})
export class ReportComponent implements OnInit, AfterViewInit {
  @ViewChild('distribution') distributionChart!: ElementRef<HTMLDivElement>;
  @ViewChild('tierHistory') tierHistoryChart!: ElementRef<HTMLDivElement>;
  @ViewChild('netChange') netChangeChart!: ElementRef<HTMLDivElement>;

  private readonly destroyRef = inject(DestroyRef);
  private resizeObserver?: ResizeObserver;
  private viewReady = false;
  private trendData?: ReportChartResponse;

  memberChurn = 0;
  membershipDistributionData: ReportChartResponse[] = [];
  selectedDistribution?: ReportChartResponse;

  constructor(private readonly reportService: ReportService) {}

  ngOnInit(): void {
    forkJoin({
      churnResponse: this.reportService.getMemberChurn(),
      reportsResponse: this.reportService.getReportsCharts(),
    })
      .pipe(takeUntilDestroyed(this.destroyRef))
      .subscribe(({ churnResponse, reportsResponse }) => {
        this.memberChurn = churnResponse.churn;
        this.trendData = reportsResponse[0];
        this.membershipDistributionData = reportsResponse.slice(1);
        this.selectedDistribution = this.membershipDistributionData.at(-1);
        this.renderCharts();
      });
  }

  ngAfterViewInit(): void {
    this.viewReady = true;
    this.resizeObserver = new ResizeObserver(() => this.renderCharts());
    this.resizeObserver.observe(this.distributionChart.nativeElement);
    this.resizeObserver.observe(this.tierHistoryChart.nativeElement);
    this.resizeObserver.observe(this.netChangeChart.nativeElement);
    this.destroyRef.onDestroy(() => this.resizeObserver?.disconnect());
    this.renderCharts();
  }

  get currentMemberCount(): number {
    return this.toChartData(this.trendData).at(-1)?.value ?? 0;
  }

  get selectedMonth(): string {
    return this.selectedDistribution?.options.title ?? 'Select a month';
  }

  updateMembershipDistributionChart(chartData: ReportChartResponse): void {
    this.selectedDistribution = chartData;
    this.renderDistributionChart();
  }

  private renderCharts(): void {
    if (!this.viewReady) return;
    this.renderDistributionChart();
    this.renderTierHistoryChart();
    this.renderNetChangeChart();
  }

  private renderDistributionChart(): void {
    if (!this.selectedDistribution || !this.viewReady) return;

    const host = this.distributionChart.nativeElement;
    const data = this.toChartData(this.selectedDistribution).filter(
      (datum) => datum.value > 0
    );
    const width = Math.max(host.clientWidth, 280);
    const height = width < 520 ? 380 : 340;
    const radius = Math.min(width * 0.28, height * 0.38);

    select(host).selectAll('*').remove();
    const svg = select(host)
      .append('svg')
      .attr('viewBox', `0 0 ${width} ${height}`)
      .attr('role', 'img')
      .attr(
        'aria-label',
        `${this.selectedMonth} membership distribution: ${data
          .map((datum) => `${datum.label} ${datum.value}`)
          .join(', ')}`
      );

    if (!data.length) {
      this.renderEmptyState(svg, width, height);
      return;
    }

    const total = sum(data, (datum) => datum.value);
    const centerX = width < 520 ? width / 2 : width * 0.32;
    const centerY = width < 520 ? height * 0.42 : height / 2;
    const color = (label: string): string =>
      TIER_COLORS.get(this.toTierKey(label)) ?? '#7c8796';
    const pieLayout = pie<ChartDatum>()
      .sort(null)
      .value((datum) => datum.value);
    const arcShape = arc<PieArcDatum<ChartDatum>>()
      .innerRadius(radius * 0.64)
      .outerRadius(radius);
    const group = svg
      .append('g')
      .attr('transform', `translate(${centerX},${centerY})`);
    const tooltip = this.createTooltip(host);

    group
      .selectAll('path')
      .data(pieLayout(data))
      .join('path')
      .attr('class', 'donut-slice')
      .attr('fill', (datum) => color(datum.data.label))
      .attr('d', arcShape)
      .attr('tabindex', 0)
      .attr('aria-label', (datum) =>
        this.tooltipText(datum.data, total)
      )
      .on('pointerenter focus', (event, datum) => {
        select(event.currentTarget).classed('active', true);
        this.showTooltip(tooltip, this.tooltipText(datum.data, total), event);
      })
      .on('pointermove', (event) =>
        this.positionTooltip(tooltip, event)
      )
      .on('pointerleave blur', (event) => {
        select(event.currentTarget).classed('active', false);
        tooltip.style('opacity', 0);
      });

    group
      .append('text')
      .attr('class', 'donut-total')
      .attr('text-anchor', 'middle')
      .attr('dy', '-0.05em')
      .text(total);
    group
      .append('text')
      .attr('class', 'donut-caption')
      .attr('text-anchor', 'middle')
      .attr('dy', '1.35em')
      .text('members');

    const legendX = width < 520 ? width * 0.12 : width * 0.61;
    const legendY = width < 520 ? height * 0.78 : centerY - data.length * 18;
    const legend = svg
      .append('g')
      .attr('class', 'chart-legend')
      .attr('transform', `translate(${legendX},${legendY})`);
    const items = legend
      .selectAll('g')
      .data(data)
      .join('g')
      .attr('transform', (_, index) => `translate(0,${index * 36})`);
    items
      .append('circle')
      .attr('r', 6)
      .attr('fill', (datum) => color(datum.label));
    items
      .append('text')
      .attr('x', 16)
      .attr('dy', '0.35em')
      .text(
        (datum) =>
          `${datum.label}  ${datum.value} (${format('.0%')(
            datum.value / total
          )})`
      );
  }

  private renderTierHistoryChart(): void {
    if (!this.viewReady) return;

    const host = this.tierHistoryChart.nativeElement;
    const data = this.toTierHistoryData();
    const width = Math.max(host.clientWidth, 280);
    const height = width < 520 ? 380 : 430;
    const margin = { top: 62, right: width < 520 ? 28 : 112, bottom: 52, left: 52 };
    const innerWidth = width - margin.left - margin.right;
    const innerHeight = height - margin.top - margin.bottom;

    select(host).selectAll('*').remove();
    const svg = select(host)
      .append('svg')
      .attr('viewBox', `0 0 ${width} ${height}`)
      .attr('role', 'img')
      .attr(
        'aria-label',
        `Membership tier history: ${data
          .map(
            (month) =>
              `${month.month}, total ${this.tierTotal(month)}, ${TIER_KEYS.map(
                (tier) => `${tier} ${month[tier]}`
              ).join(', ')}`
          )
          .join('; ')}`
      );

    if (!data.length) {
      this.renderEmptyState(svg, width, height);
      return;
    }

    const x = scaleBand<string>()
      .domain(data.map((datum) => datum.month))
      .range([0, innerWidth])
      .padding(0.18);
    const stackedData = stack<TierHistoryDatum>().keys(TIER_KEYS)(data);
    const totals = data.map((datum) => ({
      month: datum.month,
      value: this.tierTotal(datum),
    }));
    const maxTotal = max(totals, (datum) => datum.value) ?? 1;
    const y = scaleLinear()
      .domain([0, maxTotal])
      .nice()
      .range([innerHeight, 0]);
    const visibleTicks = this.visibleMonthTicks(
      data.map((datum) => datum.month),
      width,
      width < 520 ? 4 : 7
    );
    const group = svg
      .append('g')
      .attr('transform', `translate(${margin.left},${margin.top})`);

    group
      .append('g')
      .attr('class', 'grid-lines')
      .call(axisLeft(y).ticks(5).tickSize(-innerWidth).tickFormat(() => ''));
    group
      .append('g')
      .attr('class', 'chart-axis')
      .attr('transform', `translate(0,${innerHeight})`)
      .call(axisBottom(x).tickValues(visibleTicks));
    group
      .append('g')
      .attr('class', 'chart-axis')
      .call(axisLeft(y).ticks(5).tickFormat(format('d')));

    const segments = stackedData.flatMap((series) =>
      series.map((point) => ({
        month: point.data.month,
        tier: series.key as TierKey,
        value: point.data[series.key as TierKey],
        y0: point[0],
        y1: point[1],
      })).filter((point) => point.value > 0)
    );
    const tooltip = this.createTooltip(host);

    group
      .selectAll('.tier-segment')
      .data(segments)
      .join('rect')
      .attr('class', 'tier-segment')
      .attr('x', (datum) => x(datum.month) ?? 0)
      .attr('y', (datum) => y(datum.y1))
      .attr('width', x.bandwidth())
      .attr('height', (datum) => Math.max(0, y(datum.y0) - y(datum.y1)))
      .attr('fill', (datum) => TIER_COLORS.get(datum.tier) ?? '#7c8796')
      .attr('tabindex', 0)
      .attr(
        'aria-label',
        (datum) => `${datum.month}, ${datum.tier}: ${datum.value} members`
      )
      .on('pointerenter focus', (event, datum) => {
        select(event.currentTarget).classed('active', true);
        this.showTooltip(
          tooltip,
          `${datum.month} · ${datum.tier}: ${datum.value} members`,
          event
        );
      })
      .on('pointermove', (event) => this.positionTooltip(tooltip, event))
      .on('pointerleave blur', (event) => {
        select(event.currentTarget).classed('active', false);
        tooltip.style('opacity', 0);
      });

    const totalLine = line<ChartDatum>()
      .x((datum) => (x(datum.label) ?? 0) + x.bandwidth() / 2)
      .y((datum) => y(datum.value))
      .curve(curveMonotoneX);
    const totalData = totals.map((datum) => ({
      label: datum.month,
      value: datum.value,
    }));
    group
      .append('path')
      .datum(totalData)
      .attr('class', 'total-line')
      .attr('d', totalLine);

    const latest = totalData.at(-1);
    if (latest) {
      const latestX = (x(latest.label) ?? 0) + x.bandwidth() / 2;
      const latestY = y(latest.value);
      group
        .append('circle')
        .attr('class', 'latest-point')
        .attr('cx', latestX)
        .attr('cy', latestY)
        .attr('r', 4);
      group
        .append('text')
        .attr('class', 'latest-label')
        .attr('x', Math.min(latestX + 10, innerWidth + 8))
        .attr('y', latestY - 10)
        .text(`${latest.label}: ${latest.value}`);
    }

    const legend = svg
      .append('g')
      .attr('class', 'chart-legend tier-legend')
      .attr('transform', `translate(${margin.left},20)`);
    const legendItemWidth = Math.min(130, innerWidth / TIER_KEYS.length);
    const legendItems = legend
      .selectAll('g')
      .data(TIER_KEYS)
      .join('g')
      .attr('transform', (_, index) => `translate(${index * legendItemWidth},0)`);
    legendItems
      .append('rect')
      .attr('width', 12)
      .attr('height', 12)
      .attr('rx', 3)
      .attr('fill', (tier) => TIER_COLORS.get(tier) ?? '#7c8796');
    legendItems
      .append('text')
      .attr('x', 18)
      .attr('y', 10)
      .text((tier) => tier);
  }

  private renderNetChangeChart(): void {
    if (!this.trendData || !this.viewReady) return;

    const host = this.netChangeChart.nativeElement;
    const totals = this.toChartData(this.trendData);
    const data: NetChangeDatum[] = totals.slice(1).map((datum, index) => ({
      month: datum.label,
      value: datum.value - totals[index].value,
    }));
    const width = Math.max(host.clientWidth, 280);
    const height = width < 520 ? 320 : 360;
    const margin = { top: 22, right: 24, bottom: 52, left: 52 };
    const innerWidth = width - margin.left - margin.right;
    const innerHeight = height - margin.top - margin.bottom;

    select(host).selectAll('*').remove();
    const svg = select(host)
      .append('svg')
      .attr('viewBox', `0 0 ${width} ${height}`)
      .attr('role', 'img')
      .attr(
        'aria-label',
        `Monthly net membership change: ${data
          .map((datum) => `${datum.month} ${datum.value >= 0 ? '+' : ''}${datum.value}`)
          .join(', ')}`
      );

    if (!data.length) {
      this.renderEmptyState(svg, width, height);
      return;
    }

    const x = scaleBand<string>()
      .domain(data.map((datum) => datum.month))
      .range([0, innerWidth])
      .padding(0.22);
    const largestChange =
      max(data, (datum) => Math.abs(datum.value)) ?? 1;
    const y = scaleLinear()
      .domain([-Math.max(1, largestChange), Math.max(1, largestChange)])
      .nice()
      .range([innerHeight, 0]);
    const zeroY = y(0);
    const visibleTicks = this.visibleMonthTicks(
      data.map((datum) => datum.month),
      width,
      width < 520 ? 4 : 7
    );
    const group = svg
      .append('g')
      .attr('transform', `translate(${margin.left},${margin.top})`);

    group
      .append('g')
      .attr('class', 'grid-lines')
      .call(axisLeft(y).ticks(5).tickSize(-innerWidth).tickFormat(() => ''));
    group
      .append('g')
      .attr('class', 'chart-axis')
      .attr('transform', `translate(0,${innerHeight})`)
      .call(axisBottom(x).tickValues(visibleTicks));
    group
      .append('g')
      .attr('class', 'chart-axis')
      .call(
        axisLeft(y)
          .ticks(5)
          .tickFormat((value) => `${Number(value) > 0 ? '+' : ''}${format('d')(value)}`)
      );
    group
      .append('line')
      .attr('class', 'zero-line')
      .attr('x1', 0)
      .attr('x2', innerWidth)
      .attr('y1', zeroY)
      .attr('y2', zeroY);

    const tooltip = this.createTooltip(host);
    group
      .selectAll('.net-change-bar')
      .data(data)
      .join('rect')
      .attr('class', (datum) =>
        `net-change-bar ${datum.value > 0 ? 'gain' : datum.value < 0 ? 'loss' : 'flat'}`
      )
      .attr('x', (datum) => x(datum.month) ?? 0)
      .attr('y', (datum) => (datum.value >= 0 ? y(datum.value) : zeroY))
      .attr('width', x.bandwidth())
      .attr('height', (datum) => Math.max(2, Math.abs(y(datum.value) - zeroY)))
      .attr('rx', 3)
      .attr('tabindex', 0)
      .attr(
        'aria-label',
        (datum) =>
          `${datum.month}: ${datum.value >= 0 ? '+' : ''}${datum.value} members`
      )
      .on('pointerenter focus', (event, datum) => {
        select(event.currentTarget).classed('active', true);
        this.showTooltip(
          tooltip,
          `${datum.month}: ${datum.value >= 0 ? '+' : ''}${datum.value} members`,
          event
        );
      })
      .on('pointermove', (event) => this.positionTooltip(tooltip, event))
      .on('pointerleave blur', (event) => {
        select(event.currentTarget).classed('active', false);
        tooltip.style('opacity', 0);
      });

    const latest = data.at(-1);
    if (latest) {
      const latestX = (x(latest.month) ?? 0) + x.bandwidth() / 2;
      const latestY = latest.value >= 0 ? y(latest.value) : y(latest.value) + 18;
      group
        .append('text')
        .attr('class', 'latest-label')
        .attr('x', latestX)
        .attr('y', latestY - 8)
        .attr('text-anchor', 'middle')
        .text(`${latest.value >= 0 ? '+' : ''}${latest.value}`);
    }
  }

  private toTierHistoryData(): TierHistoryDatum[] {
    return this.membershipDistributionData.map((chart) => {
      const month: TierHistoryDatum = {
        month: chart.options.title,
        Classic: 0,
        Standard: 0,
        Premium: 0,
        Credited: 0,
      };

      for (const datum of this.toChartData(chart)) {
        month[this.toTierKey(datum.label)] = datum.value;
      }

      return month;
    });
  }

  private toTierKey(label: string): TierKey {
    return (
      TIER_KEYS.find((tier) => tier.toLowerCase() === label.toLowerCase()) ??
      'Credited'
    );
  }

  private toChartData(chart?: ReportChartResponse): ChartDatum[] {
    return (chart?.rows ?? [])
      .map((row) => ({ label: String(row[0]), value: Number(row[1]) }))
      .filter((datum) => Number.isFinite(datum.value));
  }

  private tierTotal(datum: TierHistoryDatum): number {
    return sum(TIER_KEYS, (tier) => datum[tier]);
  }

  private visibleMonthTicks(labels: string[], width: number, target: number): string[] {
    const tickEvery = Math.max(1, Math.ceil(labels.length / target));
    const visible = labels.filter(
      (_, index) => index % tickEvery === 0 || index === labels.length - 1
    );

    if (width > 520 && visible.length > 1) {
      const previous = labels.indexOf(visible[visible.length - 2]);
      if (labels.length - 1 - previous < Math.ceil(tickEvery / 2)) {
        visible.splice(visible.length - 2, 1);
      }
    }

    return visible;
  }

  private tooltipText(datum: ChartDatum, total: number): string {
    return `${datum.label}: ${datum.value} members · ${format('.1%')(
      datum.value / total
    )}`;
  }

  private createTooltip(host: HTMLDivElement): Selection<HTMLDivElement, unknown, null, undefined> {
    return select(host).append('div').attr('class', 'chart-tooltip');
  }

  private showTooltip(
    tooltip: Selection<HTMLDivElement, unknown, null, undefined>,
    content: string,
    event: PointerEvent | FocusEvent
  ): void {
    tooltip.text(content).style('opacity', 1);
    this.positionTooltip(tooltip, event);
  }

  private positionTooltip(
    tooltip: Selection<HTMLDivElement, unknown, null, undefined>,
    event: PointerEvent | FocusEvent
  ): void {
    const host = (tooltip.node()?.parentElement as HTMLElement).getBoundingClientRect();
    if ('clientX' in event && event.clientX) {
      tooltip
        .style('left', `${event.clientX - host.left + 12}px`)
        .style('top', `${event.clientY - host.top - 18}px`);
    } else if (event.target instanceof SVGElement) {
      const target = event.target.getBoundingClientRect();
      tooltip
        .style('left', `${target.left - host.left + target.width / 2}px`)
        .style('top', `${target.top - host.top - 8}px`);
    }
  }

  private renderEmptyState(
    svg: Selection<SVGSVGElement, unknown, null, undefined>,
    width: number,
    height: number
  ): void {
    svg
      .append('text')
      .attr('class', 'empty-chart')
      .attr('x', width / 2)
      .attr('y', height / 2)
      .attr('text-anchor', 'middle')
      .text('No report data available');
  }
}
