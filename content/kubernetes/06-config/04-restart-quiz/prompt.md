# Picking Up a Changed Variable

Your brewd pods still have `DB_PORT=5432` in their environment, but the ConfigMap now says 6432. The database team wants brewd on the new port by lunch, **without dropping any orders**.

Remember from chapter 3: a Deployment only rolls out when its **pod template** changes, and changing a ConfigMap doesn't change the template.

**Which command gets every brewd pod onto the new value, safely?**
