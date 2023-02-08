SELECT SUM(r.referral_commission) AS previous_cycle_commission
FROM   affiliate_manager_db.referral_table r
WHERE  r.affiliate_id = ?
  AND r.booking_time >= ?
  AND r.booking_time <= ?
GROUP  BY affiliate_id;