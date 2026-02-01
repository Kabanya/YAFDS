import { useMemo } from 'react'
import { getVisibleStages, getStageState } from '../utils/orderStages'

/**
 * OrderStageStepper - A beautiful glassmorphic stepper component for visualizing order stages
 *
 * @param {Object} props
 * @param {Object} props.order - The order object with at least { id, status }
 * @param {string} props.role - The current role ('customer' | 'courier' | 'restaurant')
 * @param {boolean} props.compact - Whether to show compact version (default: true)
 * @param {boolean} props.showTimestamps - Whether to show timestamps (default: false)
 * @param {string} props.orientation - Layout orientation ('horizontal' | 'vertical', default: 'horizontal')
 */
export default function OrderStageStepper({
  order,
  role,
  compact = true,
  showTimestamps = false,
  orientation = 'horizontal'
}) {
  // Get visible stages based on role and current status
  const visibleStages = useMemo(() => {
    return getVisibleStages(role, order?.status)
  }, [role, order?.status])

  // Early return if no stages or order
  if (!order?.status || visibleStages.length === 0) {
    return null
  }

  const isVertical = orientation === 'vertical'

  return (
    <div className={`stage-stepper ${isVertical ? 'stage-stepper-vertical' : 'stage-stepper-horizontal'}`}>
      {visibleStages.map((stage, index) => {
        const stageState = getStageState(stage.key, order.status, role)
        const isLast = index === visibleStages.length - 1

        return (
          <div key={stage.key} className="stage-wrapper">
            <div className={`stage-node stage-node-${stageState} ${compact ? 'stage-node-compact' : ''}`}>
              <div className="stage-icon">
                {stage.icon}
              </div>
              {!compact && (
                <div className="stage-label">
                  {stage.label}
                </div>
              )}
            </div>

            {/* Connector line (not shown on last stage or in compact mode) */}
            {!isLast && !isVertical && (
              <div className={`stage-connector ${stageState === 'completed' ? 'stage-connector-active' : ''}`} />
            )}
          </div>
        )
      })}
    </div>
  )
}
