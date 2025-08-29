-- records: حذف القيد القديم إذا موجود
ALTER TABLE records
DROP CONSTRAINT IF EXISTS records_request_id_fkey;

-- records: إضافة القيد الجديد
ALTER TABLE records
ADD CONSTRAINT records_request_id_fkey
FOREIGN KEY (request_id) REFERENCES medications(id);


-- information: حذف القيد القديم إذا موجود
ALTER TABLE information
DROP CONSTRAINT IF EXISTS information_request_id_fkey;

-- information: إضافة القيد الجديد
ALTER TABLE information
ADD CONSTRAINT information_request_id_fkey
FOREIGN KEY (request_id) REFERENCES medications(id);
