use std::collections::HashMap;
use std::ffi::OsString;

type Env = HashMap<OsString, OsString>;

pub fn empty() -> Env {
    HashMap::new()
}

pub fn inherit() -> Env {
    ::std::env::vars_os().collect()
}

// Diffing

pub enum Change {
    Add(OsString),
    Remove(OsString),
    Replace(OsString, OsString),
}

pub type Diff = HashMap<OsString, Change>;

pub fn diff(a: Env, b: Env) -> Diff {
    let mut diff = Diff::new();
    for (key, val) in &a {
        diff.insert(key.to_os_string(), Change::Add(val.to_os_string()));
    }
    for (key, val) in &b {
        let change = match diff.get(key) {
            None => Change::Remove(val.to_os_string()),
            Some(&Change::Add(ref old_val)) => {
                if old_val == val {
                    // Hack, re-use Change::Add to remove from the diff if it's the same k=v
                    Change::Add(old_val.to_os_string())
                } else {
                    Change::Replace(old_val.to_os_string(), val.to_os_string())
                }
            }
            Some(_) => panic!("bug"),
        };
        match change {
            // Hack, see above
            Change::Add(_) => diff.remove(key),
            _ => diff.insert(key.to_os_string(), change),
        };
    }
    diff
}
